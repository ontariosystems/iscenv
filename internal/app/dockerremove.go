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
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

func DockerRemove(instanceName string) (string, error) {
	instanceName = strings.ToLower(instanceName)
	instance := GetInstances().Find(instanceName)

	if instance != nil {
		err := DockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: instance.ID, RemoveVolumes: true, Force: true})
		return instance.ID, err
	} else {
		return "", fmt.Errorf("No such instance, name: %s\n", instanceName)
	}
}
