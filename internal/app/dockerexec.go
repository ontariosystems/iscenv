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
	"os"

	"github.com/ontariosystems/iscenv/iscenv"

	"github.com/fsouza/go-dockerclient"
)

type DockerExecError struct {
	ExitCode int
}

func (dee DockerExecError) Error() string {
	return fmt.Sprintf("Failing exit code returned from exec, exit code: %d", dee.ExitCode)
}

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
