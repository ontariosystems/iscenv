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
	"io"
	"os"
	"os/exec"
	"syscall"
)

// TODO: Get rid of ssh replace with Docker exec
func ManagedSSHFn(sshbin string, args []string) error {
	cmd := exec.Command(sshbin, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	// TODO: Consider removing the whole "quiet" feature, if not pass in a parameter to this
	//if !startQuiet {
	go io.Copy(os.Stdout, stdoutPipe)
	//}

	return cmd.Wait()
}

/// This replaces the current process with the ssh process.
func ProcessReplacingSSHFn(sshbin string, args []string) error {
	sargs := append([]string{"ssh"}, args...)
	return syscall.Exec(sshbin, sargs, []string{"TERM=xterm"})
}
