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

	"github.com/ontariosystems/iscenv/v3/internal/app"
	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing ISC product containers",
	Long:  "List the existing ISC product instance containers.  This is searching for any containers whose names begin with \"iscenv-\".  So, it is possible to confuse this command.",
	Run:   list,
}

func init() {
	rootCmd.AddCommand(listCmd)

	flags.AddFlag(listCmd, "no-trunc", false, "Don't truncate output")
	flags.AddFlagP(listCmd, "quiet", "q", false, "Only display numeric IDs")
}

func list(cmd *cobra.Command, _ []string) {
	instances := app.GetInstances()
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	if !flags.GetBool(cmd, "quiet") {
		fmt.Fprintln(w, "CONTAINER ID\tVERSION\tCREATED\tSTATUS\tSUPERSERVER\tWEB\tNAME")
	}

	for _, instance := range instances {
		id := instance.ID
		if !flags.GetBool(cmd, "no-trunc") {
			id = id[:12]
		}
		if !flags.GetBool(cmd, "quiet") {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t%s\n",
				id,
				instance.Version,
				time.Unix(instance.Created, 0).Format(time.RFC3339),
				instance.Status,
				instance.Ports.SuperServer,
				instance.Ports.Web,
				instance.Name)
		} else {
			fmt.Println(id)
		}

	}
	w.Flush()
}
