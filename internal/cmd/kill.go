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

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill INSTANCE [INSTANCE...]",
	Short: "Kill an ISC product instance",
	Long:  "Kill a running ISC product instance container without attempting a safe shutdown",
	Run:   kill,
}

func init() {
	rootCmd.AddCommand(killCmd)

	addMultiInstanceFlags(killCmd, "kill")
}

func kill(_ *cobra.Command, args []string) {
	instances := multiInstanceFlags.getInstances(args)
	for _, instanceName := range instances {
		instance, ilog := app.FindInstanceAndLogger(instanceName)
		if instance == nil {
			ilog.Error(app.ErrNoSuchInstance)
			continue
		}

		if err := app.DockerClient.KillContainer(docker.KillContainerOptions{ID: instance.ID}); err != nil {
			app.ErrorLogger(ilog, err).Fatal("Failed to kill instance")
		}

		ilog.Info("Killed instance")
	}
}
