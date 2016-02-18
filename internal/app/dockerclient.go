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
	docker "github.com/fsouza/go-dockerclient"
)

var DockerClient *docker.Client

const (
	DockerSocket = "unix:///var/run/docker.sock"
)

func init() {
	dc, err := docker.NewClient(DockerSocket)
	if err != nil {
		Fatalf("Could not open Docker client, socket: %s\n", DockerSocket)
	}

	DockerClient = dc
}
