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
	"fmt"
	"strconv"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/iscenv"
)

// GetInstances will return a list of available ISC product Instances.
// This function will panic in the case that it cannot list the containers from Docker.  Normally, these kinds of abrupt exits should be avoided outside of the actual executable command portions of the code.  However, in this case the extreme nature of the failure, the rarity of occurrence, and the large amount of error handling that would need to be added throughout the code without the exit seems to justify its existence.
func GetInstances() iscenv.ISCInstances {
	containers, err := DockerClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		// This should never happen on a normal working system.  As such, we will panic rather than forcing an error return that
		// is pointless under normal circumstances.
		ErrorLogger(nil, err).Panic("Failed to list containers")
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
				ErrorLogger(nil, err).WithField("containerID", apiContainer.ID).Warn("Failed to inspect container")
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
			ssPort, httpPort, hcPort, err := getPorts(container)
			if err != nil {
				ErrorLogger(ilog, err).Error("Failed to determine ports")
			}

			for intPort, bindings := range container.HostConfig.PortBindings {
				bp, err := GetDockerBindingPort(bindings)
				// This should *never* happen but we should still handle it
				if err != nil {
					ErrorLogger(ilog, err).WithField("internalPort", intPort).Error("Failed to parse port binding")
					continue
				}

				switch intPort {
				case DockerPort(ssPort):
					instance.Ports.SuperServer = bp
				case DockerPort(httpPort):
					instance.Ports.Web = bp
				case DockerPort(hcPort):
					instance.Ports.HealthCheck = bp
				}
			}

			instances = append(instances, instance)
		}
	}

	return instances
}

func getPorts(container *docker.Container) (ssPort, httpPort, hcPort iscenv.ContainerPort, err error) {
	for _, env := range container.Config.Env {
		setPort(iscenv.EnvInternalSS, env, &ssPort)
		setPort(iscenv.EnvInternalWeb, env, &httpPort)
		setPort(iscenv.EnvInternalHC, env, &hcPort)
	}

	missing := []string{}
	addMissing(iscenv.EnvInternalSS, ssPort, &missing)
	addMissing(iscenv.EnvInternalWeb, httpPort, &missing)
	addMissing(iscenv.EnvInternalHC, hcPort, &missing)

	if len(missing) > 0 {
		err = fmt.Errorf("Missing port environment variables: %s", strings.Join(missing, ","))
	}

	return
}

func setPort(envVar string, envString string, port *iscenv.ContainerPort) (err error) {
	prefix := envVar + "="
	if strings.HasPrefix(envString, prefix) {
		val := strings.TrimPrefix(envString, prefix)
		iport, err := strconv.Atoi(val)
		*port = iscenv.ContainerPort(iport)
		return err
	}

	return
}

func addMissing(portVar string, port iscenv.ContainerPort, missing *[]string) {
	if port == 0 {
		*missing = append(*missing, portVar)
	}
}
