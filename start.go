/*
Copyright 2015 Ontario Systems

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

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var startRemove bool
var startVersion string
var startPortOffset int64
var startQuiet bool
var volumesFrom []string
var containerLinks []string
var startCacheKeyUrl string

var startCommand = &cobra.Command{
	Use:   "start INSTANCE [INSTANCE...]",
	Short: "Start an ISC product container",
	Long:  "Create or start one or more ISC product containers with the supplied options",
}

func init() {
	startCommand.Run = start
	startCommand.Flags().BoolVarP(&startRemove, "rm", "", false, "Remove existing instance before starting.  By default, this switch will preserve port settings when recreating the instance.")
	startCommand.Flags().StringVarP(&startVersion, "version", "v", "", "The version of ISC product to start.  By default this will find the latest version on your system.")
	startCommand.Flags().StringSliceVar(&containerLinks, "link", nil, "Add link to another container.  They should be in the format 'iscenv-{iscenvname}', 'iscenv-{iscenvname}:{alias}' or '{containername}:{alias}'")
	startCommand.Flags().Int64VarP(&startPortOffset, "port-offset", "p", -1, "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
	startCommand.Flags().BoolVarP(&startQuiet, "quiet", "q", false, "Only display numeric IDs")
	startCommand.Flags().StringSliceVar(&volumesFrom, "volumes-from", nil, "Mount volumes from the specified container(s)")
	startCommand.Flags().StringVarP(&startCacheKeyUrl, "license-key-url", "k", "", "Download the cache.key file from the provided location rather than the default Statler URL")
	addMultiInstanceFlags(startCommand, "start")
}

func start(_ *cobra.Command, args []string) {
	if startVersion == "" {
		startVersion = getVersions().latest().version
	}

	instances := multiInstanceFlags.getInstances(args)
	// This loop is somewhat inefficient (with the multiple getInstances())  I doubt there will ever be enough to make it a performance issue.
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		current := getInstances()

		var offset int64 = -1

		existing := current.find(instance)
		if existing != nil {
			if !startRemove {
				nqf(startQuiet, "Ensuring instance '%s' is started...\n", instance)
				// NOTE: I wish this wasn't necessary (and maybe it isn't) but it seems that the api uses a blank hostConfig instead of nil which seems to wipe out all of the settings
				dockerClient.StartContainer(existing.ID, existing.container().HostConfig)

				if startPortOffset >= 0 {
					epo := int64(existing.portOffset())
					if epo != startPortOffset {
						nqf(startQuiet, "WARNING: The port offset for an existing instance differs from the offset specified, instance: %s, existing: %d, specified: %d\n", instance, epo, startPortOffset)
					}

					// Even if we're just starting a container we need to bump the ascending port counter so the next instance doesn't collide with this one
					startPortOffset++
				}
				executePrepWithOpts(existing, []string{})
				fmt.Println(existing.ID)
				continue
			}

			offset = int64(existing.portOffset())
			rm(nil, []string{instance})
			current = getInstances() // Reset this so an instance doesn't collide with itself at the port offset check below
		}

		nqf(startQuiet, "Creating instance '%s'...\n", instance)

		if offset == -1 {
			if startPortOffset >= 0 {
				offset = startPortOffset
				startPortOffset++
			} else {
				offset = current.calculatePortOffset()
			}
			nqf(startQuiet, "Calculated port offset as %d...\n", offset)
		} else {
			nqf(startQuiet, "Reusing port offset of %d...\n", offset)
		}

		if !current.usedPortOffset(offset) {
			container := createInstance(instance, startVersion, offset, volumesFrom, containerLinks)
			existing := getInstances().find(instance)
			if existing == nil {
				fatalf("Error loading instance after creation, name: %s", instance)
			}
			executePrep(existing)
			fmt.Println(container.ID)
		} else {
			nqf(startQuiet, "ERROR: Could not create instance due to port offset conflict, instance: %s, port offset: %d\n", instance, offset)
		}
	}
}

func createInstance(instance string, version string, portOffset int64, volumesFrom []string, containerLinks []string) *docker.Container {
	container, err := dockerClient.CreateContainer(getCreateOpts(instance, version, portOffset, volumesFrom, containerLinks))
	if err != nil {
		fatalf("Could not create instance, name: %s\n%v", instance, err)
	}

	err = dockerClient.StartContainer(container.ID, getStartOpts(portOffset, volumesFrom, containerLinks))
	if err != nil {
		fatalf("Could not start created instance, name: %s\n%v", instance, err)
	}

	return container
}

func getCreateOpts(name string, version string, portOffset int64, volumesFrom []string, containerLinks []string) docker.CreateContainerOptions {
	image := REPOSITORY + ":" + version

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
			exeDir() + ":/iscenv:ro",
			homeDir() + ":" + homeDir(),
		},
		Links: containerLinks,
		PortBindings: map[docker.Port][]docker.PortBinding{
			port(INTERNAL_PORT_SSH): portBinding(EXTERNAL_PORT_SSH, portOffset),
			port(INTERNAL_PORT_SS):  portBinding(EXTERNAL_PORT_SS, portOffset),
			port(INTERNAL_PORT_WEB): portBinding(EXTERNAL_PORT_WEB, portOffset),
		},
		VolumesFrom: volumesFrom,
	}
}

func executePrep(ensInstance *ISCInstance) {
	opts := []string{
		"-u", strconv.Itoa(os.Getuid()),
		"-g", strconv.Itoa(os.Getgid()),
		"-c", hgcachePath(),
	}

	executePrepWithOpts(ensInstance, opts)
}

func executePrepWithOpts(ensInstance *ISCInstance, opts []string) {
	hostIp, err := getDocker0InterfaceIP()
	if err == nil {
		opts = append(opts, "-i", hostIp)
	} else {
		nqf(startQuiet, "WARNING: Could not find docker0's address, 'host' entry will not be added to /etc/hosts, err: %s\n", err)
	}

	fmt.Println("Waiting for SSH...")
	err = waitForPort(hostIp, ensInstance.Ports.SSH.String(), 60*time.Second)
	if err == nil {
		fmt.Println("\tSuccess!")
	} else {
		fatalf("Error while waiting for SSH, name: %s, error: %s", ensInstance.Name, err)
	}

	if startCacheKeyUrl != "" {
		opts = append(opts, "-k", startCacheKeyUrl)
	}

	baseopts := []string{
		"/iscenv/iscenv", "_prep",
	}

	sshExec(ensInstance.Name, osSshFn, append(baseopts, opts...)...)
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		fatalf("Could not determine current user, err: %s\n", err)
	}

	return usr.HomeDir
}

func exeDir() string {
	folder, err := osext.ExecutableFolder()
	if err != nil {
		fatalf("Could not determine executable folder, err: %s\n", err)
	}

	return folder
}

func hgcachePath() string {
	out, err := exec.Command("hg", "showconfig", "extensions.hg-cache").CombinedOutput()
	if err != nil {
		fatal("hg showconfig extensions.hg-cache failed")
	}

	return expandHome(filepath.Join(filepath.Dir(filepath.Dir(strings.TrimSpace(string(out)))), "cache"))
}

func expandHome(path string) string {
	if path[:1] == "~" {
		return strings.Replace(path, "~", homeDir(), 1)
	}

	return path
}

func getDockerLogs(containerId string) ([]string, error) {
	var buf bytes.Buffer
	// TODO: There's probably a better way to do this with follow and continuous reading from the stream
	err := dockerClient.Logs(docker.LogsOptions{
		Container:    containerId,
		OutputStream: &buf,
		Stdout:       true,
		Stderr:       true,
		Timestamps:   true,
		Follow:       false,
	})

	if err != nil {
		return []string{}, err
	}

	return strings.Split(buf.String(), "\n"), nil
}

func containerName(instance string) string {
	return CONTAINER_PREFIX + instance
}

func svcUpLine(name string) string {
	return fmt.Sprintf("success: %s entered RUNNING state, process has stayed up for > than 1 seconds", name)
}

func allTrue(items map[string]bool) bool {
	for _, i := range items {
		if !i {
			return false
		}
	}

	return true
}
