/*
Copyright 2015 Ontario Systems

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
	"os"

	"github.com/spf13/cobra"
)

const (
	dockerSocket    = "unix:///var/run/docker.sock"
	containerPrefix = "iscenv-"
	registry        = "quay.io"
	repository      = registry + "/ontsys/centos-ensemble"

	portInternalSSH = 22
	portExternalSSH = 22822

	portInternalSS = 56772
	portExternalSS = 56772

	portInternalWeb = 57772
	portExternalWeb = 57772
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

func execute() {
	addCommands()
	iscenvCommand.Execute()
}

func addCommands() {
	// Container management and use
	iscenvCommand.AddCommand(startCommand)
	iscenvCommand.AddCommand(stopCommand)
	iscenvCommand.AddCommand(restartCommand)
	iscenvCommand.AddCommand(killCommand)
	iscenvCommand.AddCommand(rmCommand)
	iscenvCommand.AddCommand(sshCommand)
	iscenvCommand.AddCommand(csessionCommand)
	iscenvCommand.AddCommand(listCommand)
	iscenvCommand.AddCommand(tailCommand)
	iscenvCommand.AddCommand(purgeJournalCommand)

	// Image management
	iscenvCommand.AddCommand(versionsCommand)
	iscenvCommand.AddCommand(pullCommand)

	// Bonus!
	iscenvCommand.AddCommand(apacheCommand)

	// ISCEnv information
	iscenvCommand.AddCommand(versionCommand)

	// Internal commands
	iscenvCommand.AddCommand(internalPrepCommand)
	iscenvCommand.AddCommand(internalPurgeJournalCommand)
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
