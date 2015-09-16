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
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var internalPurgeJournalCommand = &cobra.Command{
	Use:   "_purgejournal",
	Short: "internal: purge old journal files",
	Long:  "DO NOT RUN THIS OUTSIDE OF AN INSTANCE CONTAINER. deletes all isc journal files that are not the current active journal file",
}

func init() {
	internalPurgeJournalCommand.Run = internalPurgeJournal
}

var internalPurgeJournalLastFileInfo os.FileInfo = nil
var internalPurgeJournalLastFilePath string

func internalPurgeJournal(_ *cobra.Command, _ []string) {
	// verify we are running in a container
	ensureWithinContainer("_prep")

	internalPurgeJournalLastFileInfo = nil
	err := filepath.Walk("/data/journal/", findLastVisit)

	if err != nil {
		fatalf("Failed to find journal files: ", err)
	}

	err = filepath.Walk("/data/journal/", deleteNotLastVisit)

	if err != nil {
		fatalf("Failed to delete journal files: ", err)
	}
}

func findLastVisit(path string, f os.FileInfo, err error) error {

	if f.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, "cache.lck") {
		return nil
	}

	if internalPurgeJournalLastFileInfo == nil || f.ModTime().After(internalPurgeJournalLastFileInfo.ModTime()) {
		internalPurgeJournalLastFileInfo = f
		internalPurgeJournalLastFilePath = path
	}

	return nil
}

func deleteNotLastVisit(path string, f os.FileInfo, err error) error {

	if path == internalPurgeJournalLastFilePath {
		return nil
	}

	if f.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, "cache.lck") {
		return nil
	}

	// one last check to verify it is a journal file and not something unexpected
	matched, err := regexp.MatchString(`^/data/journal/[0-9]{8}\.[0-9]{3}$`, path)
	if matched {
		os.Remove(path)
		fmt.Printf(" - deleted: %s\n", path)
	}

	return nil
}
