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
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/iscenv"
)

type DockerStartOptions struct {
	Name             string
	Repository       string
	Version          string
	PortOffset       int64
	PortOffsetSearch bool
	Environment      []string
	Volumes          []string
	VolumesFrom      []string
	ContainerLinks   []string

	Recreate bool
}

func (opts *DockerStartOptions) ToCreateContainerOptions() *docker.CreateContainerOptions {
	return &docker.CreateContainerOptions{
		Name: opts.ContainerName(),
		Config: &docker.Config{
			Image:    opts.Repository + ":" + opts.Version,
			Hostname: opts.Name,
			Env:      opts.Environment,
			Volumes:  opts.InternalVolumes(),
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
		PortBindings: map[docker.Port][]docker.PortBinding{
			DockerPort(iscenv.PortInternalSS):  DockerPortBinding(iscenv.PortExternalSS, opts.PortOffset),
			DockerPort(iscenv.PortInternalWeb): DockerPortBinding(iscenv.PortExternalWeb, opts.PortOffset),
		},
		VolumesFrom: opts.VolumesFrom,
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
