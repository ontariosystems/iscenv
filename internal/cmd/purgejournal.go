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
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"

	"github.com/spf13/cobra"
)

var purgeJournalCmd = &cobra.Command{
	Use:   "purgejournal INSTANCE",
	Short: "purge old journal files",
	Long:  "deletes all isc journal files that are not the current active journal file",
	Run:   purgeJournal,
}

func init() {
	rootCmd.AddCommand(purgeJournalCmd)

	addMultiInstanceFlags(purgeJournalCmd, "purgejournal")
}

func purgeJournal(cmd *cobra.Command, args []string) {
	instances := getMultipleInstances(cmd, args)
	for _, instanceName := range instances {
		instance, ilog := app.FindInstanceAndLogger(instanceName)
		if instance == nil {
			ilog.Error(app.ErrNoSuchInstance)
			continue
		}

		if err := app.DockerExec(instance, false, iscenv.InternalISCEnvPath, "_purgejournal"); err != nil {
			app.ErrorLogger(ilog, err).Fatal("Failed to purge journals")
		}
	}
}
