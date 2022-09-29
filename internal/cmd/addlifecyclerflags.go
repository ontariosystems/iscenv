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
	"github.com/ontariosystems/iscenv/v3/internal/app"
	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	"github.com/ontariosystems/iscenv/v3/internal/plugins"
	"github.com/spf13/cobra"
)

func addLifecyclerFlagsIfNeeded(cmd *cobra.Command) error {
	if skipPluginActivation() {
		return nil
	}

	return addLifecyclerFlags(cmd)
}

// Add the flags from the available lifecycler plugins to the provided command
func addLifecyclerFlags(cmd *cobra.Command) error {
	// Logging can't have been configured yet, so we're using an empty PluginArgs
	var lcs []*plugins.ActivatedLifecycler
	defer getActivatedLifecyclers(nil, plugins.PluginArgs{}, &lcs)()

	available := make([]string, len(lcs))
	for i, lc := range lcs {
		available[i] = lc.Id
		pluginFlags, err := lc.Lifecycler.Flags()
		if err != nil {
			return app.NewPluginError(lc.Id, "Flags", lc.Path, err)
		}

		for _, pluginFlag := range pluginFlags.Flags {
			flagName := lc.Id + "-" + pluginFlag.Flag
			if pluginFlag.HasConfig {
				flags.AddConfigFlag(cmd, flagName, pluginFlag.DefaultValue, pluginFlag.Usage)
			} else {
				flags.AddFlag(cmd, flagName, pluginFlag.DefaultValue, pluginFlag.Usage)
			}
		}

	}

	return nil
}
