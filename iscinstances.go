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
	"strconv"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type ContainerPort int64

type containerPorts struct {
	SSH         ContainerPort
	SuperServer ContainerPort
	Web         ContainerPort
}

type ISCInstance struct {
	ID      string
	Name    string
	Version string
	Created int64
	Status  string
	Ports   containerPorts
}

type ISCInstances []ISCInstance

func (p ContainerPort) String() string {
	return strconv.FormatInt(int64(p), 10)
}

func (i ISCInstance) portOffset() ContainerPort {
	if i.Ports.SSH < EXTERNAL_PORT_SSH {
		fatalf("SSH Port is outside of range, instance: %s, port: %s\n", i.Name, i.Ports.SSH)
	}

	return i.Ports.SSH - EXTERNAL_PORT_SSH
}

func (i ISCInstance) container() *docker.Container {
	container, err := dockerClient.InspectContainer(i.ID)
	if err != nil {
		fatalf("Could not inspect container, instance: %s, id: %s\n", i.Name, i.ID)
	}

	return container
}

func (is ISCInstances) byPortOffsets() map[ContainerPort]ISCInstance {
	offsets := make(map[ContainerPort]ISCInstance)
	for _, i := range is {
		offsets[i.portOffset()] = i
	}

	return offsets
}

func (is ISCInstances) calculatePortOffset() int64 {
	offsets := is.byPortOffsets()

	var i ContainerPort
	for i = 0; i < 65535; i++ {
		if _, in := offsets[i]; !in {
			return int64(i)
		}
	}

	fatal("Could not determine next port offset")
	return -1
}

func (is ISCInstances) usedPortOffset(offset int64) bool {
	offsets := is.byPortOffsets()
	_, used := offsets[ContainerPort(offset)]
	return used
}

func (is ISCInstances) find(name string) *ISCInstance {
	for _, i := range is {
		if i.Name == name {
			return &i
		}
	}

	return nil
}

func (is ISCInstances) exists(name string) bool {
	return is.find(name) != nil
}

func getInstances() ISCInstances {
	containers, err := dockerClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fatalf("Could not list containers, error: %s\n", err)
	}

	instances := []ISCInstance{}
	for _, apicontainer := range containers {
		name := ""

		for _, cn := range apicontainer.Names {
			// Skip over link/name container names.  Root names will be "/{CONTAINER_PREFIX}-{name}".
			if strings.Count(cn, "/") == 1 && strings.HasPrefix(cn, "/"+CONTAINER_PREFIX) {
				name = cn
				break
			}
		}

		if name != "" {
			container, err := dockerClient.InspectContainer(apicontainer.ID)
			if err != nil {
				fatalf("Could not inspect container, id: %s, error: %s\n", apicontainer.ID, err)
			}

			var version string
			if strings.Contains(apicontainer.Image, ":") {
				version = strings.Split(apicontainer.Image, ":")[1]
			} else {
				version = "Unknown"
			}

			instance := ISCInstance{
				ID:      container.ID,
				Name:    strings.TrimPrefix(name, "/"+CONTAINER_PREFIX),
				Version: version,
				Status:  apicontainer.Status,
				Created: apicontainer.Created}

			for intPort, bindings := range container.HostConfig.PortBindings {
				switch intPort {
				case port(INTERNAL_PORT_SSH):
					instance.Ports.SSH = getBindingPort(bindings)
				case port(INTERNAL_PORT_SS):
					instance.Ports.SuperServer = getBindingPort(bindings)
				case port(INTERNAL_PORT_WEB):
					instance.Ports.Web = getBindingPort(bindings)
				}
			}

			instances = append(instances, instance)
		}
	}

	return instances
}
