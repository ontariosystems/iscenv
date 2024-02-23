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
	"github.com/ontariosystems/iscenv/v3/internal/app"
	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var csessionCmd = &cobra.Command{
	Use:    "csession INSTANCE",
	Short:  "Start csession for instance",
	Long:   "[DEPRECATED] Connect to an instance container and initiate a csession. Use `iscenv session` instead",
	Hidden: true,
	Run:    csession,
}

var sessionCmd = &cobra.Command{
	Use:   "session INSTANCE",
	Short: "Start terminal session for instance",
	Long:  "Connect to an instance container and initiate a ISC terminal session.",
	Run:   csession,
}

func init() {
	rootCmd.AddCommand(csessionCmd)
	addDockerUserFlags(csessionCmd)
	flags.AddFlag(csessionCmd, "exec", "", "Execute the following command with csession")
	flags.AddConfigFlagP(csessionCmd, "namespace", "n", "%SYS", "Use a specific starting namespace")
	flags.AddConfigFlagP(csessionCmd, "internal-instance", "i", "docker", "The name of the actual ISC product instance within the container")
	flags.AddConfigFlag(csessionCmd, "control-command", "ccontrol", "The command used to control the ISC product instance")

	rootCmd.AddCommand(sessionCmd)
	addDockerUserFlags(sessionCmd)
	flags.AddFlag(sessionCmd, "exec", "", "Execute the following command with ISC terminal session")
	flags.AddConfigFlagP(sessionCmd, "namespace", "n", "%SYS", "Use a specific starting namespace")
	flags.AddConfigFlagP(sessionCmd, "internal-instance", "i", "iris", "The name of the actual ISC product instance within the container")
	flags.AddConfigFlag(sessionCmd, "control-command", "iris", "The command used to control the ISC product instance")
}

func csession(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logAndExit(log.WithError(app.ErrSingleInstanceArg), "Invalid arguments")
	}

	instanceName := flags.GetString(cmd, "internal-instance")
	control := flags.GetString(cmd, "control-command")

	cmdArgs := []string{control, "session", instanceName}
	ns := flags.GetString(cmd, "namespace")
	if ns != "" {
		cmdArgs = append(cmdArgs, "-U")
		cmdArgs = append(cmdArgs, ns)
	}

	exec := flags.GetString(cmd, "exec")
	if exec != "" {
		cmdArgs = append(cmdArgs, exec)
	}

	instance, ilog := app.FindInstanceAndLogger(args[0])
	if instance == nil {
		logAndExit(ilog.WithError(app.ErrNoSuchInstance), "Invalid arguments")
	}

	username := flags.GetString(cmd, userFlag)
	if err := app.DockerExec(instance, true, username, cmdArgs...); err != nil {
		logAndExit(app.ErrorLogger(ilog, err), "Failed to run docker exec")
	}
}
