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

package iscenv

import (
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type ISCInstances []*ISCInstance

func (is ISCInstances) ByPortOffsets() map[ContainerPort]*ISCInstance {
	offsets := make(map[ContainerPort]*ISCInstance)
	for _, i := range is {
		offsets[i.PortOffset()] = i
	}

	return offsets
}

func (is ISCInstances) CalculatePortOffset() int64 {
	offsets := is.ByPortOffsets()

	var i ContainerPort
	for i = 0; i < 65535; i++ {
		if _, in := offsets[i]; !in {
			return int64(i)
		}
	}

	Fatal("Could not determine next port offset")
	return -1
}

func (is ISCInstances) UsedPortOffset(offset int64) bool {
	offsets := is.ByPortOffsets()
	_, used := offsets[ContainerPort(offset)]
	return used
}

func (is ISCInstances) Find(name string) *ISCInstance {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}

	return nil
}

func (is ISCInstances) Exists(name string) bool {
	return is.Find(name) != nil
}

func GetInstances() ISCInstances {
	containers, err := DockerClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		Fatalf("Could not list containers, error: %s\n", err)
	}

	instances := make([]*ISCInstance, 0)
	for _, apiContainer := range containers {
		name := ""

		for _, cn := range apiContainer.Names {
			// Skip over link/name container names.  Root names will be "/{CONTAINER_PREFIX}-{name}".
			if strings.Count(cn, "/") == 1 && strings.HasPrefix(cn, "/"+ContainerPrefix) {
				name = cn
				break
			}
		}

		if name != "" {
			container, err := DockerClient.InspectContainer(apiContainer.ID)
			if err != nil {
				Fatalf("Could not inspect container, id: %s, error: %s\n", apiContainer.ID, err)
			}

			var version string
			if strings.Contains(apiContainer.Image, ":") {
				version = strings.Split(apiContainer.Image, ":")[1]
			} else {
				version = "Unknown"
			}

			instance := &ISCInstance{
				ID:      container.ID,
				Name:    strings.TrimPrefix(name, "/"+ContainerPrefix),
				Version: version,
				Status:  apiContainer.Status,
				Created: apiContainer.Created}

			for intPort, bindings := range container.HostConfig.PortBindings {
				switch intPort {
				case DockerPort(PortInternalSS):
					instance.Ports.SuperServer = GetDockerBindingPort(bindings)
				case DockerPort(PortInternalWeb):
					instance.Ports.Web = GetDockerBindingPort(bindings)
				}
			}

			instances = append(instances, instance)
		}
	}

	return instances
}
