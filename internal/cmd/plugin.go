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

	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/internal/app"
)

func init() {
	rootCmd.AddCommand(pluginCmd)
	for ptype, plugins := range app.InternalPlugins {
		typeCmd := &cobra.Command{
			Use:    ptype,
			Short:  fmt.Sprintf("Start a %s plugin", ptype),
			Long:   fmt.Sprintf("Start a plugin of type %s.", ptype),
			Hidden: true,
		}
		pluginCmd.AddCommand(typeCmd)

		for key := range plugins {
			plugin := plugins[key]
			cmd := &cobra.Command{
				Use:    key,
				Short:  fmt.Sprintf("Start the %s plugin", ptype),
				Long:   fmt.Sprintf("Start the %s %s plugin.", key, ptype),
				Hidden: true,
				Run: func(_ *cobra.Command, _ []string) {
					plugin.Main()
				},
			}

			typeCmd.AddCommand(cmd)
		}
	}
}

var pluginCmd = &cobra.Command{
	Use:    "plugin",
	Short:  "Start an embedded plugin",
	Long:   "Start an embedded plugin.",
	Hidden: true,
}
