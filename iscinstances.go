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

type containerPort int64

type containerPorts struct {
	SSH         containerPort
	SuperServer containerPort
	Web         containerPort
}

type iscInstance struct {
	ID      string
	Name    string
	Version string
	Created int64
	Status  string
	Ports   containerPorts
}

type iscInstances []iscInstance

func (p containerPort) String() string {
	return strconv.FormatInt(int64(p), 10)
}

func (i iscInstance) portOffset() containerPort {
	if i.Ports.SSH < portExternalSSH {
		fatalf("SSH Port is outside of range, instance: %s, port: %s\n", i.Name, i.Ports.SSH)
	}

	return i.Ports.SSH - portExternalSSH
}

func (i iscInstance) container() *docker.Container {
	container, err := dockerClient.InspectContainer(i.ID)
	if err != nil {
		fatalf("Could not inspect container, instance: %s, id: %s\n", i.Name, i.ID)
	}

	return container
}

func (is iscInstances) byPortOffsets() map[containerPort]iscInstance {
	offsets := make(map[containerPort]iscInstance)
	for _, i := range is {
		offsets[i.portOffset()] = i
	}

	return offsets
}

func (is iscInstances) calculatePortOffset() int64 {
	offsets := is.byPortOffsets()

	var i containerPort
	for i = 0; i < 65535; i++ {
		if _, in := offsets[i]; !in {
			return int64(i)
		}
	}

	fatal("Could not determine next port offset")
	return -1
}

func (is iscInstances) usedPortOffset(offset int64) bool {
	offsets := is.byPortOffsets()
	_, used := offsets[containerPort(offset)]
	return used
}

func (is iscInstances) find(name string) *iscInstance {
	for _, i := range is {
		if i.Name == name {
			return &i
		}
	}

	return nil
}

func (is iscInstances) exists(name string) bool {
	return is.find(name) != nil
}

func getInstances() iscInstances {
	containers, err := dockerClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fatalf("Could not list containers, error: %s\n", err)
	}

	instances := []iscInstance{}
	for _, apicontainer := range containers {
		name := ""

		for _, cn := range apicontainer.Names {
			// Skip over link/name container names.  Root names will be "/{CONTAINER_PREFIX}-{name}".
			if strings.Count(cn, "/") == 1 && strings.HasPrefix(cn, "/"+containerPrefix) {
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

			instance := iscInstance{
				ID:      container.ID,
				Name:    strings.TrimPrefix(name, "/"+containerPrefix),
				Version: version,
				Status:  apicontainer.Status,
				Created: apicontainer.Created}

			for intPort, bindings := range container.HostConfig.PortBindings {
				switch intPort {
				case port(portInternalSSH):
					instance.Ports.SSH = getBindingPort(bindings)
				case port(portInternalSS):
					instance.Ports.SuperServer = getBindingPort(bindings)
				case port(portInternalWeb):
					instance.Ports.Web = getBindingPort(bindings)
				}
			}

			instances = append(instances, instance)
		}
	}

	return instances
}
