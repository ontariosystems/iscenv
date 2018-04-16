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
	"os"
	"path/filepath"
	"sort"

	"github.com/ontariosystems/iscenv/internal/app"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var internalPurgeJournalCmd = &cobra.Command{
	Use:   "_purgejournal",
	Short: "internal: purge old journal files",
	Long:  "deletes all isc journal files that are not the current active journal file (this command is only available within containers)",
	Run:   internalPurgeJournal,
}

func init() {
	if err := app.EnsureWithinContainer("_purgejournal"); err != nil {
		return
	}

	rootCmd.AddCommand(internalPurgeJournalCmd)
}

func internalPurgeJournal(_ *cobra.Command, _ []string) {
	journals, err := filepath.Glob("/data/journal/[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9].[0-9][0-9][0-9]")
	if err != nil {
		logAndExit(app.ErrorLogger(log.StandardLogger(), err), "Failed to list journal files")
	}

	if len(journals) < 2 {
		log.Info("No old journal files found")
		return
	}

	sort.Strings(journals)
	for _, journal := range journals[0 : len(journals)-1] {
		jlog := log.WithField("journalFile", journal)

		f, err := os.Open(journal)
		if err != nil {
			logAndExit(app.ErrorLogger(jlog, err), "Failed to open journal file")
		}

		fi, err := f.Stat()
		if err != nil {
			logAndExit(app.ErrorLogger(jlog, err), "Failed to stat journal file")
		}

		if err := os.Remove(journal); err != nil {
			logAndExit(app.ErrorLogger(jlog, err), "Failed to remove journal file")
		}

		jlog.WithField("journalSize", fi.Size()/1024/1024).Info("Deleted journal file")
	}
}
