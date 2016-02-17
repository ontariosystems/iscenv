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
	"text/tabwriter"
	"time"

	"github.com/ontariosystems/iscenv/internal/app"

	"github.com/spf13/cobra"
)

var versionsFlags = &struct {
	NoTrunc bool
	Quiet   bool
}{}

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List existing ISC product versions",
	Long:  "List the existing ISC product versions.  This finds ISC product images with tags matching the version format.",
	Run:   versions,
}

func init() {
	rootCmd.AddCommand(versionsCmd)

	versionsCmd.Flags().BoolVarP(&versionsFlags.NoTrunc, "no-trunc", "", false, "Don't truncate output")
	versionsCmd.Flags().BoolVarP(&versionsFlags.Quiet, "quiet", "q", false, "Only display numeric IDs")
}

func versions(_ *cobra.Command, _ []string) {
	versions := app.GetVersions()
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	if !versionsFlags.Quiet {
		fmt.Fprintln(w, "IMAGE ID\tVERSION\tCREATED")
	}

	for _, version := range versions {
		id := version.ID
		if !versionsFlags.NoTrunc {
			id = id[:12]
		}
		if !versionsFlags.Quiet {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				id,
				version.Version,
				time.Unix(version.Created, 0).Format(time.RFC3339))
		} else {
			fmt.Println(id)
		}

	}
	w.Flush()
}
