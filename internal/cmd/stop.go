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
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [OPTIONS] INSTANCE [INSTANCE...]",
	Short: "Stop an ISC product instance",
	Long:  "Stop a running ISC product instance, attempting a safe shutdown",
	Run:   stop,
}

func init() {
	rootCmd.AddCommand(stopCmd)

	addMultiInstanceFlags(stopCmd, "stop")
	flags.AddFlagP(stopCmd, "time", "t", uint(60), "The number of seconds to wait for the instance to stop cleanly before killing it.")
}

func stop(cmd *cobra.Command, args []string) {
	instances := getMultipleInstances(cmd, args)
	for _, instanceName := range instances {
		instance, ilog := app.FindInstanceAndLogger(instanceName)
		if instance == nil {
			ilog.Error(app.ErrNoSuchInstance)
			continue
		}

		if err := app.DockerClient.StopContainer(instance.ID, flags.GetUint(cmd, "time")); err != nil {
			app.ErrorLogger(ilog, err).Fatal("Failed to stop instance")
		}

		ilog.Info("Stopped instance")
	}
}