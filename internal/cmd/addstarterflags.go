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
	"strings"

	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
)

// Adding starter flags when doing plugin calls causes an infinite loop
func addStarterFlagsIfNotPluginCall(cmd *cobra.Command) error {
	if isPluginCall() {
		return nil
	}

	return addStarterFlags(cmd)
}

// Add the flags from the available starter plugins to the provided command
func addStarterFlags(cmd *cobra.Command) error {
	available := make([]string, 0)
	// Logging can't have been configured yet, so we're using an empty PluginArgs
	if err := activateStartersAndClose(nil, app.PluginArgs{}, func(id, path string, starter iscenv.Starter) error {
		available = append(available, id)
		pluginFlags, err := starter.Flags()
		if err != nil {
			return app.NewPluginError(id, "Flags", path, err)
		}

		for _, pluginFlag := range pluginFlags.Flags {
			flagName := id + "-" + pluginFlag.Flag
			if pluginFlag.HasConfig {
				flags.AddConfigFlag(cmd, flagName, pluginFlag.DefaultValue, pluginFlag.Usage)
			} else {
				flags.AddFlag(cmd, flagName, pluginFlag.DefaultValue, pluginFlag.Usage)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	flags.AddConfigFlag(cmd, "plugins", "", "An ordered comma-separated list of plugins you wish to activate. available plugins: "+strings.Join(available, ","))

	return nil
}
