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
	"strings"

	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/iscenv"
)

// Add the flags from the available starter plugins to the provided command
func addStarterFlags(cmd *cobra.Command, pluginsFlag *string, pluginFlags *iscenv.PluginFlags) error {
	available := make([]string, 0)
	if err := activateStartersAndClose(nil, func(id, _ string, starter iscenv.Starter) error {
		var err error
		available = append(available, id)
		startFlags.PluginFlags, err = starter.Flags()
		if err != nil {
			return fmt.Errorf("Could not retrieve plugin flags, plugin: %s, error: %s", id, err)
		}
		pluginFlags.AddFlagsToFlagSet(id, cmd.Flags())
		return nil
	}); err != nil {
		return err
	}

	cmd.Flags().StringVar(pluginsFlag, "plugins", "", "An ordered comma-separated list of plugins you wish to activate. available plugins: "+strings.Join(available, ","))

	return nil
}
