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
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ontariosystems/iscenv/internal/iscenv"

	"github.com/spf13/cobra"
)

var internalPurgeJournalCmd = &cobra.Command{
	Use:    "_purgejournal",
	Short:  "internal: purge old journal files",
	Long:   "DO NOT RUN THIS OUTSIDE OF AN INSTANCE CONTAINER. deletes all isc journal files that are not the current active journal file",
	Hidden: true,
	Run:    internalPurgeJournal,
}

func init() {
	rootCmd.AddCommand(internalPurgeJournalCmd)
}

func internalPurgeJournal(_ *cobra.Command, _ []string) {
	// verify we are running in a container
	iscenv.EnsureWithinContainer("_purgejournal")

	journals, err := filepath.Glob("/data/journal/[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9].[0-9][0-9][0-9]")
	if err != nil {
		iscenv.Fatalf("Failed to list journal files: %s", err)
	}

	if len(journals) < 2 {
		fmt.Printf("  - no old journal files found\n")
		return
	}

	sort.Strings(journals)
	for _, journal := range journals[0 : len(journals)-1] {
		f, err := os.Open(journal)
		if err != nil {
			iscenv.Fatalf("Could not open journal file %s: %s", journal, err)
		}

		fi, err := f.Stat()
		if err != nil {
			iscenv.Fatalf("Could not stat journal file %s: %s", journal, err)
		}

		if err := os.Remove(journal); err != nil {
			iscenv.Fatalf("Could not delete journal, path: %s, error: %s", journal, err)
		} else {
			fmt.Printf("  - deleted: %s (%v MB)\n", journal, fi.Size()/1024/1024)
		}
	}
}
