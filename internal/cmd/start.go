/*
Copyright 2016 Ontario Systems

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var startFlags = struct {
	Remove         bool
	Version        string
	PortOffset     int64
	Quiet          bool
	VolumesFrom    []string
	ContainerLinks []string
	CacheKeyURL    string
}{}

var startCmd = &cobra.Command{
	Use:   "start INSTANCE [INSTANCE...]",
	Short: "Start an ISC product container",
	Long:  "Create or start one or more ISC product containers with the supplied options",
	Run:   start,
}

func init() {
	// Since, we're adding flags and this has to happen in init, we're unfortunately going to have to load up and close the plugins here and in the start function, we could persist the manager globally but it's not as safe as a failure in init could concievably leave rpc servers running
	pm, err := app.NewPluginManager(iscenv.ApplicationName, iscenv.StarterKey, iscenv.StarterPlugin{})
	if err != nil {
		app.Fatalf("Could not load raw interfaces, error: %S", err)
	}
	defer pm.Close()

	if err := pm.VisitPlugins(func(id string, raw interface{}) error {
		starter := raw.(iscenv.Starter)
		flags, err := starter.Flags()
		if err != nil {
			return fmt.Errorf("Could not retrieve plugin flags, plugin: %s, error: %s", id, err)
		}
		flags.AddFlagsToFlagSet(id, startCmd.Flags())
		return nil
	}); err != nil {
		app.Fatalf("Failed to retrieve plugin flags, error: %s", err)
	}

	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&startFlags.Remove, "rm", "", false, "Remove existing instance before starting.  By default, this switch will preserve port settings when recreating the instance.")
	startCmd.Flags().StringVarP(&startFlags.Version, "version", "v", "", "The version of ISC product to start.  By default this will find the latest version on your system.")
	startCmd.Flags().StringSliceVar(&startFlags.ContainerLinks, "link", nil, "Add link to another container.  They should be in the format 'iscenv-{iscenvname}', 'iscenv-{iscenvname}:{alias}' or '{containername}:{alias}'")
	startCmd.Flags().Int64VarP(&startFlags.PortOffset, "port-offset", "p", -1, "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
	startCmd.Flags().BoolVarP(&startFlags.Quiet, "quiet", "q", false, "Only display numeric IDs")
	startCmd.Flags().StringSliceVar(&startFlags.VolumesFrom, "volumes-from", nil, "Mount volumes from the specified container(s)")
	startCmd.Flags().StringVarP(&startFlags.CacheKeyURL, "license-key-url", "k", "", "Download the cache.key file from the provided location rather than the default Statler URL")
	addMultiInstanceFlags(startCmd, "start")
}

// TODO: There is *way* too much logic in this command.  Move it into libraries in app
func start(_ *cobra.Command, args []string) {
	if startFlags.Version == "" {
		startFlags.Version = app.GetVersions().Latest().Version
	}

	instances := multiInstanceFlags.getInstances(args)
	// This loop is somewhat inefficient (with the multiple app.GetInstances())  I doubt there will ever be enough to make it a performance issue.
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		current := app.GetInstances()

		var offset int64 = -1
		var err error

		// TODO: Look through this logic, it seems weird...
		existing := current.Find(instance)
		if existing != nil {
			if !startFlags.Remove {
				app.Nqf(startFlags.Quiet, "Ensuring instance '%s' is started...\n", instance)
				// NOTE: I wish this wasn't necessary (and maybe it isn't) but it seems that the api uses a blank hostConfig instead of nil which seems to wipe out all of the settings
				existingContainer := app.GetContainerForInstance(existing)
				app.DockerClient.StartContainer(existing.ID, existingContainer.HostConfig)

				if startFlags.PortOffset >= 0 {
					epo, err := existing.PortOffset()
					if err != nil {
						app.Fatalf("Error determining port offset, instance: %s, error: %s", existing.Name, err)
					}

					if epo != startFlags.PortOffset {
						app.Nqf(startFlags.Quiet, "WARNING: The port offset for an existing instance differs from the offset specified, instance: %s, existing: %d, specified: %d\n", instance, epo, startFlags.PortOffset)
					}

					// Even if we're just starting a container we need to bump the ascending port counter so the next instance doesn't collide with this one
					startFlags.PortOffset++
				}
				executePrepWithOpts(existing, []string{})
				fmt.Println(existing.ID)
				continue
			}

			offset, err = existing.PortOffset()
			if err != nil {
				app.Fatalf("Error determining port offset, instance: %s, error: %s", existing.Name, err)
			}
			rm(nil, []string{instance})
			current = app.GetInstances() // Reset this so an instance doesn't collide with itself at the port offset check below
		}

		app.Nqf(startFlags.Quiet, "Creating instance '%s'...\n", instance)

		if offset == -1 {
			if startFlags.PortOffset >= 0 {
				offset = startFlags.PortOffset
				startFlags.PortOffset++
			} else {
				offset, err = current.CalculatePortOffset()
				if err != nil {
					app.Fatalf("Error calculating port offset, instance: %s, error: %s", instance, err)
				}
			}
			app.Nqf(startFlags.Quiet, "Calculated port offset as %d...\n", offset)
		} else {
			app.Nqf(startFlags.Quiet, "Reusing port offset of %d...\n", offset)
		}

		upo, err := current.UsedPortOffset(offset)
		if err != nil {
			app.Fatalf("Error checking port offset usage, instance: %s, error: %s", instance, err)
		}

		if !upo {
			container := createInstance(instance, startFlags.Version, offset, startFlags.VolumesFrom, startFlags.ContainerLinks)
			existing := app.GetInstances().Find(instance)
			if existing == nil {
				app.Fatalf("Error loading instance after creation, name: %s", instance)
			}
			executePrep(existing)
			fmt.Println(container.ID)
		} else {
			app.Nqf(startFlags.Quiet, "ERROR: Could not create instance due to port offset conflict, instance: %s, port offset: %d\n", instance, offset)
		}
	}
}

func createInstance(instance string, version string, portOffset int64, volumesFrom []string, containerLinks []string) *docker.Container {
	container, err := app.DockerClient.CreateContainer(getCreateOpts(instance, version, portOffset, startFlags.VolumesFrom, containerLinks))
	if err != nil {
		app.Fatalf("Could not create instance, name: %s\n%s", instance, err)
	}

	err = app.DockerClient.StartContainer(container.ID, getStartOpts(portOffset, volumesFrom, containerLinks))
	if err != nil {
		app.Fatalf("Could not start created instance, name: %s\n%s", instance, err)
	}

	return container
}

func getCreateOpts(name string, version string, portOffset int64, volumesFrom []string, containerLinks []string) docker.CreateContainerOptions {
	image := iscenv.Repository + ":" + version

	home := homeDir()
	config := docker.Config{
		Image:    image,
		Hostname: name,
		Env:      []string{"HOST_HOME=" + home},
		Volumes: map[string]struct{}{
			// TODO: See if we can drop the data volume
			"/data":             struct{}{},
			"/var/log/ensemble": struct{}{},
			"/iscenv":           struct{}{},
			home:                struct{}{},
		}}

	opts := docker.CreateContainerOptions{
		Name:       containerName(name),
		Config:     &config,
		HostConfig: getStartOpts(portOffset, volumesFrom, containerLinks),
	}

	return opts
}

func getStartOpts(portOffset int64, volumesFrom []string, containerLinks []string) *docker.HostConfig {
	return &docker.HostConfig{
		Privileged: true,
		Binds: []string{
			fmt.Sprintf("%s:%s:ro", exe(), iscenv.InternalISCEnvPath),
			fmt.Sprintf("%s:%s", homeDir(), homeDir()),
		},
		Links: containerLinks,
		PortBindings: map[docker.Port][]docker.PortBinding{
			app.DockerPort(iscenv.PortInternalSS):  app.DockerPortBinding(iscenv.PortExternalSS, portOffset),
			app.DockerPort(iscenv.PortInternalWeb): app.DockerPortBinding(iscenv.PortExternalWeb, portOffset),
		},
		VolumesFrom: volumesFrom,
	}
}

func executePrep(ensInstance *iscenv.ISCInstance) {
	opts := []string{
		"-u", strconv.Itoa(os.Getuid()),
		"-g", strconv.Itoa(os.Getgid()),
		"-c", hgcachePath(),
	}

	executePrepWithOpts(ensInstance, opts)
}

func executePrepWithOpts(ensInstance *iscenv.ISCInstance, opts []string) {
	hostIP, err := app.GetDocker0InterfaceIP()
	if err == nil {
		opts = append(opts, "-i", hostIP)
	} else {
		app.Nqf(startFlags.Quiet, "WARNING: Could not find docker0's address, 'host' entry will not be added to /etc/hosts, err: %s\n", err)
	}

	if startFlags.CacheKeyURL != "" {
		opts = append(opts, "-k", startFlags.CacheKeyURL)
	}

	baseopts := []string{
		iscenv.InternalISCEnvPath, "_prep",
	}

	if err := app.DockerExec(ensInstance.Name, false, append(baseopts, opts...)...); err != nil {
		app.Fatalf("Could not run prep in newly started container, instance: %s, error: %s\n", ensInstance.Name, err)
	}
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		app.Fatalf("Could not determine current user, err: %s\n", err)
	}

	return usr.HomeDir
}

func exe() string {
	exe, err := osext.Executable()
	if err != nil {
		app.Fatalf("Could not determine executable path, err: %s\n", err)
	}

	return exe
}

func hgcachePath() string {
	out, err := exec.Command("hg", "showconfig", "extensions.hg-cache").CombinedOutput()
	if err != nil {
		app.Fatal("hg showconfig extensions.hg-cache failed")
	}

	return expandHome(filepath.Join(filepath.Dir(filepath.Dir(strings.TrimSpace(string(out)))), "cache"))
}

func expandHome(path string) string {
	if path[:1] == "~" {
		return strings.Replace(path, "~", homeDir(), 1)
	}

	return path
}

func containerName(instance string) string {
	return iscenv.ContainerPrefix + instance
}
