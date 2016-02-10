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
	"archive/tar"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
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

func DockerExec(instanceName string, interactive bool, commandAndArgs ...string) error {
	instance := GetInstances().Find(strings.ToLower(instanceName))
	if instance == nil {
		return fmt.Errorf("Could not find instance, name: %s", instanceName)
	}

	createOpts := docker.CreateExecOptions{
		Container:    instance.ID,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          interactive,
		Cmd:          commandAndArgs,
	}

	exec, err := DockerClient.CreateExec(createOpts)
	if err != nil {
		return err
	}

	startOpts := docker.StartExecOptions{
		Tty:         interactive,
		RawTerminal: interactive,
	}

	var stdin *rawTTYStdin
	if interactive {
		stdin, err = NewRawTTYStdin()
		if err != nil {
			return err
		}

		startOpts.InputStream = stdin
	}

	startOpts.OutputStream = os.Stdout
	startOpts.ErrorStream = os.Stderr
	cw, err := DockerClient.StartExecNonBlocking(exec.ID, startOpts)
	if err != nil {
		return err
	}

	if cw == nil {
		return nil
	}

	if stdin != nil {
		stdin.MonitorTTYSize(func(height, width int) {
			DockerClient.ResizeExecTTY(exec.ID, height, width)
		})
	}

	return cw.Wait()
}

func DockerCopy(instanceName, instancePath, localPath string) error {
	instance := GetInstances().Find(strings.ToLower(instanceName))
	if instance == nil {
		return fmt.Errorf("Could not find instance, name: %s", instanceName)
	}

	r, w := io.Pipe()

	go func() {
		DockerClient.DownloadFromContainer(instance.ID, docker.DownloadFromContainerOptions{
			Path:         instancePath,
			OutputStream: w,
		})
	}()

	t := tar.NewReader(r)
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(localPath, header.Name)
		info := header.FileInfo()
		fmt.Println(path, info.Name())
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, t)
		if err != nil {
			return err
		}
	}

	return nil
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
