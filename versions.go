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
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var versionsNoTrunc bool
var versionsQuiet bool

var versionsCommand = &cobra.Command{
	Use:   "versions",
	Short: "List existing ISC product versions",
	Long:  "List the existing ISC product versions.  This finds ISC product images with tags matching the version format.",
}

func init() {
	versionsCommand.Run = versions
	versionsCommand.Flags().BoolVarP(&versionsNoTrunc, "no-trunc", "", false, "Don't truncate output")
	versionsCommand.Flags().BoolVarP(&versionsQuiet, "quiet", "q", false, "Only display numeric IDs")
}

func versions(_ *cobra.Command, _ []string) {
	versions := getVersions()
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	if !versionsQuiet {
		fmt.Fprintln(w, "IMAGE ID\tVERSION\tCREATED")
	}

	for _, version := range versions {
		id := version.id
		if !versionsNoTrunc {
			id = id[:12]
		}
		if !versionsQuiet {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				id,
				version.version,
				time.Unix(version.created, 0).Format(time.RFC3339))
		} else {
			fmt.Println(id)
		}

	}
	w.Flush()
}
