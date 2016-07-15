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
	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/internal/app"

	log "github.com/Sirupsen/logrus"
)

var (
	defaultExecCommand = []string{"/bin/bash", "-l"}
)

var execCmd = &cobra.Command{
	Use:   "exec INSTANCE -- [COMMAND] [ARGS...]",
	Short: "Connect to an instance",
	Long:  "Connect to an instance container with docker exec.",
	Run:   dockerExec,
}

func init() {
	rootCmd.AddCommand(execCmd)
}

func dockerExec(_ *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal(app.ErrSingleInstanceArg)
	}

	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	} else {
		cmdArgs = defaultExecCommand
	}

	instance, ilog := app.FindInstanceAndLogger(args[0])
	if instance == nil {
		ilog.Fatal(app.ErrNoSuchInstance)
	}

	if err := app.DockerExec(instance, true, cmdArgs...); err != nil {
		app.ErrorLogger(ilog, err).Fatal("Failed to run docker exec")
	}
}
