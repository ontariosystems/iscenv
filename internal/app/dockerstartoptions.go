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

package app

import (
	"strconv"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/v3/iscenv"
	log "github.com/sirupsen/logrus"
)

// DockerStartOptions holds information used in starting a container
type DockerStartOptions struct {
	// The name of the instance
	Name string

	// The full name of the container, if this is populated, it will not be calculated
	// TODO: This is fairly hacky and should be reworked in a future refactoring phase
	FullName string

	// The image repository from which the container will be created
	Repository string

	// The version of the image to use
	Version string

	// The port by which the external ports will be offset (or the starting offset for search if searching)
	PortOffset int64

	// Search for the next available port offset?
	PortOffsetSearch bool

	// Disable the port offset checks, only do this if you are creating additional supplementary containers that use an existing containers offset to offset a different port/set of ports
	// TODO: This is fairly hacky and should be reworked in a future refactoring phase
	DisablePortOffsetConflictCheck bool

	// The entrypoint for the container
	Entrypoint []string

	// The command for the container
	Command []string

	// Environment variables in standard docker format (ENV=VALUE)
	Environment []string

	// Volumes provided in the standard host:container:mode format
	Volumes []string

	// Copies files provided in the format host:container into the container before it starts
	Copies []string

	// The names of containers from which volumes will be used
	VolumesFrom []string

	// Containers to which this container should be linked
	ContainerLinks []string

	// Port mappings in standard <IP>:host:container format
	Ports []string

	// Labels to add to the container
	Labels map[string]string

	// Should the container be deleted if it already exists?
	Recreate bool

	// User to run the container as
	Username string
}

// ToCreateContainerOptions transforms a DockerStartOptions into CreateContainerOptions
func (opts *DockerStartOptions) ToCreateContainerOptions() *docker.CreateContainerOptions {
	createOpts := &docker.CreateContainerOptions{
		Name: opts.ContainerName(),
		Config: &docker.Config{
			Image:        opts.Repository + ":" + opts.Version,
			Hostname:     opts.Name,
			Env:          opts.Environment,
			Volumes:      opts.InternalVolumes(),
			ExposedPorts: opts.ToExposedPorts(),
			Entrypoint:   opts.Entrypoint,
			Cmd:          opts.Command,
			Labels:       opts.Labels,
			User:         opts.Username,
		},
		HostConfig: opts.ToHostConfig(),
	}

	return createOpts
}

// ToHostConfig transforms a DockerStartOptions into a HostConfig
func (opts *DockerStartOptions) ToHostConfig() *docker.HostConfig {

	return &docker.HostConfig{
		// TODO: Try turning this off or better still allow it to be activated with a plugin or better even again allow the appropriate capabilities to be set with a plugin
		Privileged:   true,
		Binds:        opts.VolumeBinds(),
		Links:        opts.ContainerLinks,
		CgroupnsMode: "host",

		// Plugin
		PortBindings: opts.ToDockerPortBindings(),
		VolumesFrom:  opts.VolumesFrom,
	}
}

// InternalVolumes returns a map of the volumes internal to the container extracted from the DockerStartOptions list of volumes
func (opts *DockerStartOptions) InternalVolumes() map[string]struct{} {
	volumes := make(map[string]struct{})
	for _, volume := range opts.Volumes {
		s := strings.Split(volume, ":")
		if len(s) == 1 {
			volumes[s[0]] = struct{}{}
		} else {
			volumes[s[1]] = struct{}{}
		}
	}
	return volumes
}

// VolumeBinds returns a slice of volumes that should be bound from the DockerStartOptions
func (opts *DockerStartOptions) VolumeBinds() []string {
	volumes := make([]string, 0)
	for _, volume := range opts.Volumes {
		if strings.Contains(volume, ":") {
			volumes = append(volumes, volume)
		}
	}

	return volumes
}

// ContainerName returns the name that should be used for the container that is being managed
func (opts *DockerStartOptions) ContainerName() string {
	if opts.FullName != "" {
		return opts.FullName
	}
	return iscenv.ContainerPrefix + opts.Name
}

// ToExposedPorts returns a map of ports that are exposed by the contain
func (opts *DockerStartOptions) ToExposedPorts() map[docker.Port]struct{} {
	ports := make(map[docker.Port]struct{})
	for port := range opts.ToDockerPortBindings() {
		ports[port] = struct{}{}
	}

	return ports
}

// ToDockerPortBindings returns a map of Ports to a slice of PortBindings
func (opts *DockerStartOptions) ToDockerPortBindings() map[docker.Port][]docker.PortBinding {
	bindings := make(map[docker.Port][]docker.PortBinding)

	if opts.Ports != nil {
		for _, bindString := range opts.Ports {
			s := strings.Split(bindString, ":")
			var hostIP, hostPort, containerPort string
			switch len(s) {
			case 2:
				hostIP = ""
				hostPort = s[0]
				containerPort = s[1]
			case 3:
				hostIP = s[0]
				hostPort = s[1]
				containerPort = s[2]
			default:
				log.WithField("portString", bindString).Warn("Single port mappings are not supported")
			}

			if strings.HasPrefix(hostPort, "+") {
				hostPort = strings.TrimPrefix(hostPort, "+")
				i, err := strconv.ParseInt(hostPort, 10, 64)
				if err != nil {
					log.WithField("port", hostPort).Warn("Could not parse host port")
					continue
				}
				hostPort = strconv.FormatInt(i+opts.PortOffset, 10)
			}

			cp := docker.Port(containerPort + "/tcp")
			if _, ok := bindings[cp]; !ok {
				bindings[cp] = make([]docker.PortBinding, 0)
			}

			bindings[cp] = append(bindings[cp], docker.PortBinding{HostIP: hostIP, HostPort: hostPort})
		}
	}
	return bindings
}
