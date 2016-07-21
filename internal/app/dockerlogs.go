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
// DockerLogs will retrieve logs from a docker container
// If follow is true, the logs will be followed
// since is the unix time from which logs should be retrieved
// tail is the string value to pass to docker to determine the number of lines to tail. It should be either "all" or a string representation of the number of lines.
// Will return any error encountered
func DockerLogs(instance *iscenv.ISCInstance, since int64, tail string, follow bool, outputStream io.Writer) error {
	opts := docker.LogsOptions{
		Container:    instance.ID,
		OutputStream: outputStream,
		Since:        since,
		Tail:         tail,
		Follow:       follow,
		Stdout:       true,
		Stderr:       true,
	}

	return DockerClient.Logs(opts)
}
