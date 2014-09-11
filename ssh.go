/*
Copyright 2014 Ontario Systems

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

package main

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type sshExecFn func(string, []string) error

var sshExecCommand string

var sshCommand = &cobra.Command{
	Use:   "ssh INSTANCE",
	Short: "SSH to an instance",
	Long:  "Connect to an instance with SSH using private key auth.",
}

func init() {
	sshCommand.Run = ssh
	sshCommand.Flags().StringVarP(&sshExecCommand, "command", "c", "", "Execute an SSH command directly")
}

func ssh(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		sshArgs := []string{}
		if sshExecCommand != "" {
			sshArgs = append(sshArgs, sshExecCommand)
		}
		sshExec(args[0], nil, sshArgs...)
	} else {
		fatal("Must provide an instance")
	}
}

func sshExec(instance string, sshfn sshExecFn, sshArgs ...string) {
	instance = strings.ToLower(instance)
	existing := getInstances().find(instance)
	if existing != nil {
		key, err := ioutil.TempFile("", "iscenv-key")
		if err != nil {
			fatalf("Could not create ssh key file, error: %s\n", err)
		}

		defer key.Close()
		defer os.Remove(key.Name())

		key.WriteString(SSH_KEY)
		key.Close()

		sshbin, err := exec.LookPath("ssh")
		if err != nil {
			fatalf("Could not find ssh binary on path, error: %s\n", err)
		}

		args := append([]string{
			"root@localhost",
			"-p", existing.ports.ssh.String(),
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "StrictHostKeyChecking=no",
			"-o", "LogLevel=error",
			"-i", key.Name()},
			sshArgs...)

		if sshfn == nil {
			sshfn = syscallSshFn
		}
		err = sshfn(sshbin, args)

		if err != nil {
			fatalf("ssh to instance failed, instance: %s, error: %s\n", instance, err)
		}
	} else {
		fatalf("No such instance, name: %s\n", instance)
	}
}

func syscallSshFn(sshbin string, args []string) error {
	sargs := append([]string{"ssh"}, args...)
	return syscall.Exec(sshbin, sargs, []string{})
}
