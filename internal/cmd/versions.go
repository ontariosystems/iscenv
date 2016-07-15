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

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
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
	pm := getVersionerPM()

	rootCmd.AddCommand(versionsCmd)

	flags.AddFlag(versionsCmd, "no-trunc", false, "Don't truncate output")
	flags.AddFlagP(versionsCmd, "quiet", "q", false, "Only display numeric IDs")
	flags.AddConfigFlag(versionsCmd, "plugins", "", `An ordered comma-separated list of plugins you wish to activate.  The "local" versions plugin will always be active as as the baseline. available plugins: `+strings.Join(pm.AvailablePlugins(), ","))
}

func versions(cmd *cobra.Command, _ []string) {
	ensureImage()

	// Only debug or fatal logging so we don't corrupt the table output
	plugins := strings.Split(flags.GetString(cmd, "plugins"), ",")
	image := flags.GetString(rootCmd, "image")
	versions, err := getVersions(image, plugins)
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve versions")
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

// Acquire all of the versions for the provided image using the appropriate plugin stack
func getVersions(image string, plugins []string) (iscenv.ISCVersions, error) {
	pm := getVersionerPM()
	defer pm.Close()

	// Get the baseline set of versions that are considered "local"
	var versions iscenv.ISCVersions

	// No need for error handling as we'll always log fatal within the loop in the event of an error
	log.WithField("plugin", baseVersionPlugin).Debug("Executing default version plugin")
	if err := pm.ActivatePlugins([]string{baseVersionPlugin}, func(id, path string, raw interface{}) error {
		var err error
		plog := app.PluginLogger(id, "Versions", path)
		versioner := raw.(iscenv.Versioner)

		plog.Debug("Retrieving versions")
		versions, err = versioner.Versions(image)
		if err != nil {
			plog.Fatal("Failed to load versions from plugin")
		}

		plog.WithField("count", len(versions)).Debug("Retrieved versions")
		return nil
	}); err != nil {
		log.WithError(err).Error("Execution of default version plugin failed")
	}

	// No need for error handling as we'll always log fatal within the loop in the event of an error
	log.Debugf("Executing %d additional version plugin(s)", len(plugins))
	if err := pm.ActivatePlugins(plugins, func(id, path string, raw interface{}) error {
		// Local was added to the plugins list which makes no sense but isn't worthy of an error (and we don't want to log because it will corrupt the table output of versions)
		if strings.EqualFold(id, baseVersionPlugin) {
			log.WithField("plugin", baseVersionPlugin).Debug("Skipping default version plugin (it was already executed)")
			return nil
		}

		plog := app.PluginLogger(id, "Versions", path)
		versioner := raw.(iscenv.Versioner)

		plog.Debug("Retrieving versions")
		plugVers, err := versioner.Versions(image)
		if err != nil {
			plog.Fatal("Failed to load versions from plugin")
		}

		plog.WithField("count", len(plugVers)).Debug("Retrieved versions")
		// TODO: Once some more plugins are implemented this will need to be changed to show what is local, what is remote, what is both but stale, etc.
		for _, version := range plugVers {
			versions.AddIfMissing(version)
		}

		return nil
	}); err != nil {
		log.WithError(err).Error("Execution of additional version plugins failed")
	}

	versions.Sort()
	return versions, nil
}

func getVersionerPM() *app.PluginManager {
	pm, err := app.NewPluginManager(
		iscenv.ApplicationName,
		iscenv.VersionerKey,
		iscenv.VersionerPlugin{},
		app.PluginArgs{
			LogLevel: flags.GetString(rootCmd, "log-level"),
			LogJSON:  flags.GetBool(rootCmd, "log-json"),
		},
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to create plugin manager")
	}
	return pm
}
