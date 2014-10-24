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
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	DOCKER_SOCKET    = "unix:///var/run/docker.sock"
	CONTAINER_PREFIX = "iscenv-"
	REGISTRY         = "quay.io"
	REPOSITORY       = REGISTRY + "/ontsys/centos-ensemble"

	INTERNAL_PORT_SSH = 22
	EXTERNAL_PORT_SSH = 22822

	INTERNAL_PORT_SS = 56772
	EXTERNAL_PORT_SS = 56772

	EXTERNAL_PORT_WEB = 57772
	INTERNAL_PORT_WEB = 57772
)

var verbose bool

var iscenvCommand = &cobra.Command{
	Use:   "iscenv",
	Short: "Manage Docker-based ISC product environments",
	Long:  "This tool allows the creation and management of Docker-based ISC product Environments.",
}

func init() {
	iscenvCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Verbose output")
}

func Execute() {
	AddCommands()
	iscenvCommand.Execute()
}

func AddCommands() {
	// Container management and use
	iscenvCommand.AddCommand(startCommand)
	iscenvCommand.AddCommand(stopCommand)
	iscenvCommand.AddCommand(killCommand)
	iscenvCommand.AddCommand(rmCommand)
	iscenvCommand.AddCommand(sshCommand)
	iscenvCommand.AddCommand(csessionCommand)
	iscenvCommand.AddCommand(listCommand)
	iscenvCommand.AddCommand(tailCommand)

	// Run in container
	iscenvCommand.AddCommand(prepCommand)

	// Image management
	iscenvCommand.AddCommand(versionsCommand)
	iscenvCommand.AddCommand(pullCommand)

	// ISCEnv information
	iscenvCommand.AddCommand(versionCommand)
}

func nq(quiet bool, a ...interface{}) {
	if !quiet {
		fmt.Println(a...)
	}
}

func nqf(quiet bool, format string, a ...interface{}) {
	if !quiet {
		fmt.Printf(format, a...)
	}
}

func fatal(a ...interface{}) {
	fmt.Println(a...)
	os.Exit(1)
}

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
