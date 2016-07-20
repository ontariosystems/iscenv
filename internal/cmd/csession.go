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
	"log"

	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
)

var csessionCmd = &cobra.Command{
	Use:   "csession INSTANCE",
	Short: "Start csession for instance",
	Long:  "Connect to an instance container and initiate a csession.",
	Run:   csession,
}

func init() {
	rootCmd.AddCommand(csessionCmd)
	flags.AddConfigFlagP(csessionCmd, "namespace", "n", "%SYS", "Use a specific starting namespace")
	flags.AddFlag(csessionCmd, "exec", "", "Execute the following command with csession")
}

func csession(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal(app.ErrSingleInstanceArg)
	}

	cmdArgs := []string{"csession", "docker"}
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
		ilog.Fatal(app.ErrNoSuchInstance)
	}

	if err := app.DockerExec(instance, true, cmdArgs...); err != nil {
		app.ErrorLogger(ilog, err).Fatal("Failed to run docker exec")
	}
}
