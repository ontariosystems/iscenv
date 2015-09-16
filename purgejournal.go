/*
Copyright 2015 Ontario Systems

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

package main

import (
	"github.com/spf13/cobra"
)

var purgeJournalCommand = &cobra.Command{
	Use:   "purgejournal INSTANCE",
	Short: "purge old journal files",
	Long:  "deletes all isc journal files that are not the current active journal file",
}

func init() {
	purgeJournalCommand.Run = purgeJournal
}

func purgeJournal(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		sshExec(args[0], nil, "/iscenv/iscenv", "_purgejournal")
	} else {
		fatal("Must provide an instance")
	}
}
