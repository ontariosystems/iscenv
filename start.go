/*
Copyright 2014 Ontario Systems

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
	"bitbucket.org/kardianos/osext"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var startRemove bool
var startVersion string
var startPortOffset int64
var startQuiet bool

var startCommand = &cobra.Command{
	Use:   "start INSTANCE [INSTANCE...]",
	Short: "Start an ISC product container",
	Long:  "Create or start one or more ISC product containers with the supplied options",
}

func init() {
	startCommand.Run = start
	startCommand.Flags().BoolVarP(&startRemove, "rm", "", false, "Remove existing instance before starting.  By default, this switch will preserve port settings when recreating the instance.")
	startCommand.Flags().StringVarP(&startVersion, "version", "v", "", "The version of ISC product to start.  By default this will find the latest version on your system.")
	startCommand.Flags().Int64VarP(&startPortOffset, "port-offset", "p", -1, "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
	startCommand.Flags().BoolVarP(&startQuiet, "quiet", "q", false, "Only display numeric IDs")
}

func start(_ *cobra.Command, args []string) {
	if startVersion == "" {
		startVersion = getVersions().latest().version
	}

	// This loop is somewhat inefficient (with the multiple getInstances())  I doubt there will ever be enough to make it a performance issue.
	for _, arg := range args {
		instance := strings.ToLower(arg)
		current := getInstances()

		var offset int64 = -1

		existing := current.find(instance)
		if existing != nil {
			if !startRemove {
				nqf(startQuiet, "Ensuring instance '%s' is started...\n", instance)
				// NOTE: I wish this wasn't necessary (and maybe it isn't) but it seems that the api uses a blank hostConfig instead of nil which seems to wipe out all of the settings
				dockerClient.StartContainer(existing.id, existing.container().HostConfig)

				if startPortOffset >= 0 {
					epo := int64(existing.portOffset())
					if epo != startPortOffset {
						nqf(startQuiet, "WARNING: The port offset for an existing instance differs from the offset specified, instance: %s, existing: %d, specified: %d\n", instance, epo, startPortOffset)
					}

					// Even if we're just starting a container we need to bump the ascending port counter so the next instance doesn't collide with this one
					startPortOffset++
				}
				fmt.Println(existing.id)
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
			container := createInstance(instance, startVersion, offset)
			executePrep(instance)
			fmt.Println(container.ID)
		} else {
			nqf(startQuiet, "ERROR: Could not create instance due to port offset conflict, instance: %s, port offset: %d\n", instance, offset)
		}
	}
}

func createInstance(instance string, version string, portOffset int64) *docker.Container {
	container, err := dockerClient.CreateContainer(getCreateOpts(instance, version, portOffset))
	if err != nil {
		fatalf("Could not create instance, name: %s\n", instance)
	}

	err = dockerClient.StartContainer(container.ID, getStartOpts(portOffset))
	if err != nil {
		fatalf("Could not start created instance, name: %s\n", instance)
	}

	return container
}

func getCreateOpts(name string, version string, portOffset int64) docker.CreateContainerOptions {
	image := REPOSITORY + ":" + version

	home := homeDir()
	config := docker.Config{
		Image:    image,
		Hostname: name,
		Env:      []string{"HOST_HOME=" + home},
		Volumes: map[string]struct{}{
			"/data":   struct{}{},
			"/iscenv": struct{}{},
			home:      struct{}{},
		}}

	opts := docker.CreateContainerOptions{
		Name:   CONTAINER_PREFIX + name,
		Config: &config}

	return opts
}

func getStartOpts(portOffset int64) *docker.HostConfig {
	return &docker.HostConfig{
		Privileged: true,
		Binds: []string{
			exeDir() + ":/iscenv:ro",
			homeDir() + ":" + homeDir(),
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			port(INTERNAL_PORT_SSH): portBinding(EXTERNAL_PORT_SSH, portOffset),
			port(INTERNAL_PORT_SS):  portBinding(EXTERNAL_PORT_SS, portOffset),
			port(INTERNAL_PORT_WEB): portBinding(EXTERNAL_PORT_WEB, portOffset),
		},
	}
}

func executePrep(instance string) {
	time.Sleep(5000 * time.Millisecond) //TODO: This should be something better than a sleep
	opts := []string{
		"/iscenv/iscenv", "prep",
		"-u", strconv.Itoa(os.Getuid()),
		"-g", strconv.Itoa(os.Getgid()),
		"-h", hgcachePath(),
	}

	hostIp, err := getDocker0InterfaceIP()
	if err == nil {
		opts = append(opts, "-i", hostIp)
	} else {
		nqf(startQuiet, "WARNING: Could not find docker0's address, 'host' entry will not be added to /etc/hosts, err: %s\n", err)
	}

	sshExec(instance, osSshFn, opts...)
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

	return filepath.Join(filepath.Dir(filepath.Dir(strings.TrimSpace(string(out)))), "cache")
}

func osSshFn(sshbin string, args []string) error {
	cmd := exec.Command(sshbin, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	if !startQuiet {
		go io.Copy(os.Stdout, stdoutPipe)
	}

	return cmd.Wait()
}
