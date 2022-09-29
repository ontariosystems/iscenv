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
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/v3/iscenv"
)

// DockerExecError is used for representing a failed docker exec
type DockerExecError struct {
	ExitCode int
}

// Error will return an error string for a DockerExecError
func (dee DockerExecError) Error() string {
	return fmt.Sprintf("Failing exit code returned from exec, exit code: %d", dee.ExitCode)
}

// DockerExec performs a docker exec against the provided instance
func DockerExec(instance *iscenv.ISCInstance, interactive bool, commandAndArgs ...string) error {
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
	} else {
		startOpts.InputStream = os.Stdin
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
			if err := DockerClient.ResizeExecTTY(exec.ID, height, width); err != nil {
				log.WithError(err).WithFields(log.Fields{"instance-id": exec.ID, "height": height, "width": width}).Error("cloud not resize docker tty")
			}
		})
	}

	if err := cw.Wait(); err != nil {
		return err
	}

	ei, err := DockerClient.InspectExec(exec.ID)
	if err != nil {
		return err
	}

	if ei.ExitCode > 0 {
		return DockerExecError{ExitCode: ei.ExitCode}
	}

	return nil
}
