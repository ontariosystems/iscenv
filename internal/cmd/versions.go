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
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ontariosystems/iscenv/v3/internal/app"
	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	"github.com/ontariosystems/iscenv/v3/internal/plugins"
	"github.com/ontariosystems/iscenv/v3/iscenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	baseVersionPlugin = "local"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List existing ISC product versions",
	Long:  "List the existing ISC product versions.  This finds ISC product images with tags matching the version format.",
	Run:   versions,
}

func init() {
	// We are making a pm just to get a listing of the plugins in init, we will not activate it here
	vm, err := plugins.NewVersionerManager(plugins.PluginArgs{})
	if err != nil {
		logAndExit(log.WithError(err), "Failed to load version plugin manager during init")
	}
	defer vm.Close(rootCtx)

	rootCmd.AddCommand(versionsCmd)

	flags.AddFlag(versionsCmd, "no-trunc", false, "Don't truncate output")
	flags.AddFlagP(versionsCmd, "quiet", "q", false, "Only display numeric IDs")
	flags.AddConfigFlag(versionsCmd, "plugins", "", `An ordered comma-separated list of plugins you wish to activate.  The "local" versions plugin will always be active as as the baseline. available plugins: `+strings.Join(vm.AvailablePlugins(), ","))
}

func versions(cmd *cobra.Command, _ []string) {
	ensureImage()

	// Only debug or fatal logging so we don't corrupt the table output
	image := flags.GetString(rootCmd, "image")
	versions, err := getVersions(image, getPluginsToActivate(cmd))
	if err != nil {
		logAndExit(log.WithError(err), "Failed to retrieve versions")
	}

	// No more logging at this point
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	if !flags.GetBool(cmd, "quiet") {
		fmt.Fprintln(w, "IMAGE ID\tVERSION\tCREATED\tSOURCE")
	}

	for _, version := range versions {
		id := version.ID
		if !flags.GetBool(cmd, "no-trunc") {
			id = id[:12]
		}
		if !flags.GetBool(cmd, "quiet") {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				id,
				version.Version,
				time.Unix(version.Created, 0).Format(time.RFC3339),
				version.Source,
			)
		} else {
			fmt.Println(id)
		}

	}
	w.Flush()
}

// Acquire all the versions for the provided image using the appropriate plugin stack
func getVersions(image string, verPlugins []string) (iscenv.ISCVersions, error) {
	// Get the baseline set of versions that are considered "local"
	var versions iscenv.ISCVersions
	var versioners []*plugins.ActivatedVersioner
	var err error

	log.WithField("plugin", baseVersionPlugin).Debug("Executing default version plugin")
	defer getActivatedVersioners([]string{baseVersionPlugin}, getPluginArgs(), &versioners)(rootCtx)
	if len(versioners) != 1 {
		logAndExit(log.WithField("plugin", baseVersionPlugin), "Got more than 1 plugin entry for base version plugin, this should be impossible")
	}
	base := versioners[0]
	plog := app.PluginLogger(base.Id, "Versions", base.Path)
	plog.Debug("Retrieving versions")

	versions, err = base.Versioner.Versions(image)
	if err != nil {
		logAndExit(plog.WithError(err), "Failed to load versions from plugin")
	}
	plog.WithField("count", len(versions)).Debug("Retrieved versions")

	log.WithField("count", len(verPlugins)).Debug("Executing additional version plugin(s)")
	defer getActivatedVersioners(verPlugins, getPluginArgs(), &versioners)(rootCtx)

	for _, v := range versioners {
		// Local was added to the plugins list which makes no sense but isn't worthy of an error (and we don't want to log because it will corrupt the table output of versions)
		if strings.EqualFold(v.Id, baseVersionPlugin) {
			log.WithField("plugin", baseVersionPlugin).Debug("Skipping default version plugin (it was already executed)")
			continue
		}

		plog = app.PluginLogger(v.Id, "Versions", v.Path)

		plog.Debug("Retrieving versions")
		plugVers, err := v.Versioner.Versions(image)
		if err != nil {
			logAndExit(plog, "Failed to load versions from plugin")
		}

		plog.WithField("count", len(plugVers)).Debug("Retrieved versions")
		// TODO: Once some more plugins are implemented this will need to be changed to show what is local, what is remote, what is both but stale, etc.
		for _, version := range plugVers {
			versions.AddIfMissing(version)
		}
	}
	versions.Sort()
	return versions, nil
}
