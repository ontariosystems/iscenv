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

// GetInstances will return a list of available ISC product Instances.
// This function will perform an os.Exit in the case that it cannot list the containers from Docker.  Normally, these kinds of abrupt exits should be avoided outside of the actual executable command portions of the code.  However, in this case the extreme nature of the failure, the rarity of occurrence, and the large amount of error handling that would need to be added throughout the code without the exit seems to justify its existence.
func GetInstances() iscenv.ISCInstances {
	containers, err := DockerClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		// NOTE: Fatal early exit
		ErrorLogger(nil, err).Fatal("Failed to list containers")
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
				ErrorLogger(nil, err).WithField("containerID", apiContainer.ID).Warning("Failed to inspect container")
				continue
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
				Created: apiContainer.Created,
			}

			ilog := InstanceLogger(instance)

			for intPort, bindings := range container.HostConfig.PortBindings {
				bp, err := GetDockerBindingPort(bindings)
				// This should *never* happen but we should still handle it
				if err != nil {
					ErrorLogger(ilog, err).WithField("internalPort", intPort).Error("Failed to parse port binding")
					continue
				}

				switch intPort {
				case DockerPort(iscenv.PortInternalSS):
					instance.Ports.SuperServer = bp
				case DockerPort(iscenv.PortInternalWeb):
					instance.Ports.Web = bp
				case DockerPort(iscenv.PortInternalHC):
					instance.Ports.HealthCheck = bp
				}
			}

			instances = append(instances, instance)
		}
	}

	return instances
}
