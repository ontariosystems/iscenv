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

	"github.com/ontariosystems/iscenv/iscenv"

	"github.com/fsouza/go-dockerclient"
)

func GetInstances() iscenv.ISCInstances {
	containers, err := DockerClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		Fatalf("Could not list containers, error: %s\n", err)
	}

	instances := make(iscenv.ISCInstances, 0)
	for _, apiContainer := range containers {
		name := ""

		for _, cn := range apiContainer.Names {
			// Skip over link/name container names.  Root names will be "/{CONTAINER_PREFIX}-{name}".
			if strings.Count(cn, "/") == 1 && strings.HasPrefix(cn, "/"+iscenv.ContainerPrefix) {
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

			instance := &iscenv.ISCInstance{
				ID:      container.ID,
				Name:    strings.TrimPrefix(name, "/"+iscenv.ContainerPrefix),
				Version: version,
				Status:  apiContainer.Status,
				Created: apiContainer.Created}

			for intPort, bindings := range container.HostConfig.PortBindings {
				switch intPort {
				case DockerPort(iscenv.PortInternalSS):
					instance.Ports.SuperServer = GetDockerBindingPort(bindings)
				case DockerPort(iscenv.PortInternalWeb):
					instance.Ports.Web = GetDockerBindingPort(bindings)
				}
			}

			instances = append(instances, instance)
		}
	}

	return instances
}
