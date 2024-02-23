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
	"github.com/ontariosystems/iscenv/v3/iscenv"
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

	addDockerUserFlags(purgeJournalCmd)
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

		username := flags.GetString(cmd, userFlag)
		if err := app.DockerExec(instance, false, username, iscenv.InternalISCEnvPath, "_purgejournal"); err != nil {
			logAndExit(app.ErrorLogger(ilog, err), "Failed to purge journals")
		}
	}
}
