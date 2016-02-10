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

package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ontariosystems/iscenv/internal/iscenv"

	"github.com/spf13/cobra"
)

type sshExecFn func(string, []string) error

var sshFlags = struct {
	Command string
}{}

var sshCmd = &cobra.Command{
	Use:   "ssh INSTANCE",
	Short: "SSH to an instance",
	Long:  "Connect to an instance with SSH using private key auth.",
	Run:   ssh,
}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().StringVarP(&sshFlags.Command, "command", "c", "", "Execute an SSH command directly")
}

func ssh(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		sshArgs := []string{}
		if sshFlags.Command != "" {
			sshArgs = append(sshArgs, sshFlags.Command)
		}
		sshExec(args[0], nil, sshArgs...)
	} else {
		iscenv.Fatal("Must provide an instance")
	}
}

func sshExec(instance string, sshfn sshExecFn, sshArgs ...string) {
	instance = strings.ToLower(instance)
	existing := iscenv.GetInstances().Find(instance)
	if existing != nil {
		key, err := ioutil.TempFile("", "iscenv-key")
		if err != nil {
			iscenv.Fatalf("Could not create ssh key file, error: %s\n", err)
		}

		defer key.Close()
		defer os.Remove(key.Name())

		key.WriteString(iscenv.SSHKey)
		key.Close()

		sshbin, err := exec.LookPath("ssh")
		if err != nil {
			iscenv.Fatalf("Could not find ssh binary on path, error: %s\n", err)
		}

		args := append([]string{
			"root@localhost",
			"-t",
			"-p", existing.Ports.SSH.String(),
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "StrictHostKeyChecking=no",
			"-o", "LogLevel=error",
			"-i", key.Name()},
			sshArgs...)

		if sshfn == nil {
			sshfn = iscenv.ProcessReplacingSSHFn
		}

		err = sshfn(sshbin, args)

		if err != nil {
			iscenv.Fatalf("ssh to instance failed, instance: %s, error: %s\n", instance, err)
		}
	} else {
		iscenv.Fatalf("No such instance, name: %s\n", instance)
	}
}
