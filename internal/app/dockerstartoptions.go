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

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/iscenv"
)

type DockerStartOptions struct {
	Name             string
	Repository       string
	Version          string
	PortOffset       int64
	PortOffsetSearch bool
	Entrypoint       []string
	Command          []string
	Environment      []string
	Volumes          []string
	VolumesFrom      []string
	ContainerLinks   []string
	Ports            []string

	Recreate bool
}

func (opts *DockerStartOptions) ToCreateContainerOptions() *docker.CreateContainerOptions {
	return &docker.CreateContainerOptions{
		Name: opts.ContainerName(),
		Config: &docker.Config{
			Image:        opts.Repository + ":" + opts.Version,
			Hostname:     opts.Name,
			Env:          opts.Environment,
			Volumes:      opts.InternalVolumes(),
			ExposedPorts: opts.ToExposedPorts(),
			Entrypoint:   opts.Entrypoint,
			Cmd:          opts.Command,
		},
		HostConfig: opts.ToHostConfig(),
	}
}

func (opts *DockerStartOptions) ToHostConfig() *docker.HostConfig {

	return &docker.HostConfig{
		// TODO: Try turning this off or better still allow it to be activated with a plugin or better even again allow the appropriate capabilities to be set with a plugin
		Privileged: true,
		Binds:      opts.Volumes,
		Links:      opts.ContainerLinks,
		// Plugin
		PortBindings: opts.ToDockerPortBindings(),
		VolumesFrom:  opts.VolumesFrom,
	}
}

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

func (opts *DockerStartOptions) ContainerName() string {
	return iscenv.ContainerPrefix + opts.Name
}

func (opts *DockerStartOptions) ToExposedPorts() map[docker.Port]struct{} {
	ports := make(map[docker.Port]struct{})
	for port := range opts.ToDockerPortBindings() {
		ports[port] = struct{}{}
	}

	return ports
}

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
				// TODO: log a warning?
			}

			if strings.HasPrefix(hostPort, "+") {
				strings.TrimPrefix(hostPort, "+")
				// TODO: Log a warning if search is still enabled
				i, err := strconv.ParseInt(hostPort, 10, 64)
				if err != nil {
					// TODO: log a warning
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
