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
	"fmt"
	"net"
	"strconv"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

var DockerClient *docker.Client

func init() {
	dc, err := docker.NewClient(DockerSocket)
	if err != nil {
		Fatalf("Could not open Docker client, socket: %s\n", DockerSocket)
	}

	DockerClient = dc
}

func DockerPort(port ContainerPort) docker.Port {
	return docker.Port(port.String()) + "/tcp"
}

func DockerPortBinding(port int64, portOffset int64) []docker.PortBinding {
	return []docker.PortBinding{docker.PortBinding{HostIP: "", HostPort: strconv.FormatInt(port+portOffset, 10)}}
}

// Assumes a single binding
func GetDockerBindingPort(bindings []docker.PortBinding) ContainerPort {
	port, err := strconv.ParseInt(bindings[0].HostPort, 10, 64)
	if err != nil {
		Fatalf("Could not parse port, error: %s\n", err)
	}

	return ContainerPort(port)
}

func GetDocker0InterfaceIP() (string, error) {
	i, err := net.InterfaceByName("docker0")
	if err != nil {
		return "", err
	}

	as, err := i.Addrs()
	if err != nil {
		return "", err
	}

	ip := ""
	for _, a := range as {
		ip = strings.Split(a.String(), "/")[0]
		if ip != "" {
			break
		}
	}

	if ip == "" {
		return "", fmt.Errorf("No addresses associated with docker0 device")
	}

	return ip, nil
}
