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
	"time"

	"github.com/ontariosystems/iscenv/internal/iscenv"

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

func start(_ *cobra.Command, args []string) {
	if startFlags.Version == "" {
		startFlags.Version = iscenv.GetVersions().Latest().Version
	}

	instances := multiInstanceFlags.getInstances(args)
	// This loop is somewhat inefficient (with the multiple iscenv.GetInstances())  I doubt there will ever be enough to make it a performance issue.
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		current := iscenv.GetInstances()

		var offset int64 = -1

		existing := current.Find(instance)
		if existing != nil {
			if !startFlags.Remove {
				iscenv.Nqf(startFlags.Quiet, "Ensuring instance '%s' is started...\n", instance)
				// NOTE: I wish this wasn't necessary (and maybe it isn't) but it seems that the api uses a blank hostConfig instead of nil which seems to wipe out all of the settings
				iscenv.DockerClient.StartContainer(existing.ID, existing.Container().HostConfig)

				if startFlags.PortOffset >= 0 {
					epo := int64(existing.PortOffset())
					if epo != startFlags.PortOffset {
						iscenv.Nqf(startFlags.Quiet, "WARNING: The port offset for an existing instance differs from the offset specified, instance: %s, existing: %d, specified: %d\n", instance, epo, startFlags.PortOffset)
					}

					// Even if we're just starting a container we need to bump the ascending port counter so the next instance doesn't collide with this one
					startFlags.PortOffset++
				}
				executePrepWithOpts(existing, []string{})
				fmt.Println(existing.ID)
				continue
			}

			offset = int64(existing.PortOffset())
			rm(nil, []string{instance})
			current = iscenv.GetInstances() // Reset this so an instance doesn't collide with itself at the port offset check below
		}

		iscenv.Nqf(startFlags.Quiet, "Creating instance '%s'...\n", instance)

		if offset == -1 {
			if startFlags.PortOffset >= 0 {
				offset = startFlags.PortOffset
				startFlags.PortOffset++
			} else {
				offset = current.CalculatePortOffset()
			}
			iscenv.Nqf(startFlags.Quiet, "Calculated port offset as %d...\n", offset)
		} else {
			iscenv.Nqf(startFlags.Quiet, "Reusing port offset of %d...\n", offset)
		}

		if !current.UsedPortOffset(offset) {
			container := createInstance(instance, startFlags.Version, offset, startFlags.VolumesFrom, startFlags.ContainerLinks)
			existing := iscenv.GetInstances().Find(instance)
			if existing == nil {
				iscenv.Fatalf("Error loading instance after creation, name: %s", instance)
			}
			executePrep(existing)
			fmt.Println(container.ID)
		} else {
			iscenv.Nqf(startFlags.Quiet, "ERROR: Could not create instance due to port offset conflict, instance: %s, port offset: %d\n", instance, offset)
		}
	}
}

func createInstance(instance string, version string, portOffset int64, volumesFrom []string, containerLinks []string) *docker.Container {
	container, err := iscenv.DockerClient.CreateContainer(getCreateOpts(instance, version, portOffset, startFlags.VolumesFrom, containerLinks))
	if err != nil {
		iscenv.Fatalf("Could not create instance, name: %s\n%s", instance, err)
	}

	err = iscenv.DockerClient.StartContainer(container.ID, getStartOpts(portOffset, volumesFrom, containerLinks))
	if err != nil {
		iscenv.Fatalf("Could not start created instance, name: %s\n%s", instance, err)
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
			iscenv.DockerPort(iscenv.PortInternalSSH): iscenv.DockerPortBinding(iscenv.PortExternalSSH, portOffset),
			iscenv.DockerPort(iscenv.PortInternalSS):  iscenv.DockerPortBinding(iscenv.PortExternalSS, portOffset),
			iscenv.DockerPort(iscenv.PortInternalWeb): iscenv.DockerPortBinding(iscenv.PortExternalWeb, portOffset),
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
	hostIP, err := iscenv.GetDocker0InterfaceIP()
	if err == nil {
		opts = append(opts, "-i", hostIP)
	} else {
		iscenv.Nqf(startFlags.Quiet, "WARNING: Could not find docker0's address, 'host' entry will not be added to /etc/hosts, err: %s\n", err)
	}

	fmt.Println("Waiting for SSH...")
	err = iscenv.WaitForPort(hostIP, ensInstance.Ports.SSH.String(), 60*time.Second)
	if err == nil {
		fmt.Println("\tSuccess!")
	} else {
		iscenv.Fatalf("Error while waiting for SSH, name: %s, error: %s", ensInstance.Name, err)
	}

	if startFlags.CacheKeyURL != "" {
		opts = append(opts, "-k", startFlags.CacheKeyURL)
	}

	baseopts := []string{
		iscenv.InternalISCEnvPath, "_prep",
	}

	sshExec(ensInstance.Name, iscenv.ManagedSSHFn, append(baseopts, opts...)...)
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		iscenv.Fatalf("Could not determine current user, err: %s\n", err)
	}

	return usr.HomeDir
}

func exe() string {
	exe, err := osext.Executable()
	if err != nil {
		iscenv.Fatalf("Could not determine executable path, err: %s\n", err)
	}

	return exe
}

func hgcachePath() string {
	out, err := exec.Command("hg", "showconfig", "extensions.hg-cache").CombinedOutput()
	if err != nil {
		iscenv.Fatal("hg showconfig extensions.hg-cache failed")
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
