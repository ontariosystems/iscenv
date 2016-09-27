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
	"github.com/ontariosystems/iscenv/internal/plugins"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm INSTANCE [INSTANCE...]",
	Short: "Remove an ISC product instance",
	Long:  "Forcefully remove an ISC product instance regardless of its current state",
	Run:   rm,
}

func init() {
	rootCmd.AddCommand(rmCmd)
	addMultiInstanceFlags(rmCmd, "rm")
}

func rm(cmd *cobra.Command, args []string) {
	var lcs []*plugins.ActivatedLifecycler
	defer getActivatedLifecyclers(getPluginsToActivate(rootCmd), getPluginArgs(), &lcs)()

	instances := getMultipleInstances(cmd, args)
	for _, instanceName := range instances {
		instance, ilog := app.FindInstanceAndLogger(instanceName)
		if instance == nil {
			ilog.Error(app.ErrNoSuchInstance)
			continue
		}

		ilog.WithField("count", len(lcs)).Info("Executing BeforeRemove hook(s) from plugins")
		for _, lc := range lcs {
			plog := app.PluginLogger(lc.Id, "BeforeRemove", lc.Path)
			plog.Debug("Executing BeforeRemove hook")
			if err := lc.Lifecycler.BeforeRemove(instance); err != nil {
				plog.WithError(err).Fatal("Failed to execute BeforeRemove hook")
			}
		}

		if err := app.DockerRemove(instance); err != nil {
			app.ErrorLogger(ilog, err).Fatal("Failed to remove instance")
		}

		ilog.Info("Removed instance")

	}
}
