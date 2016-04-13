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
	"io"

	"github.com/ontariosystems/iscenv/iscenv"

	"github.com/fsouza/go-dockerclient"
)

// Retrieve the logs from only this start run
func DockerLogs(instance *iscenv.ISCInstance, outputStream io.Writer) error {
	container, err := GetContainerForInstance(instance)
	if err != nil {
		return err
	}

	opts := docker.LogsOptions{
		Container:    instance.ID,
		OutputStream: outputStream,
		Since:        container.State.StartedAt.Unix(),
		Follow:       true,
		Stdout:       true,
	}

	return DockerClient.Logs(opts)
}
