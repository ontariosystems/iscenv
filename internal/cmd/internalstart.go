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
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"
)

var internalStartFlags = struct {
	Instance     string
	CControlPath string
	Plugins      string
	PluginFlags  map[string]*iscenv.PluginFlags
}{}

var internalStartCmd = &cobra.Command{
	Use:    "_start",
	Short:  "internal: start/manage ISC product ",
	Long:   "DO NOT RUN THIS OUTSIDE OF AN INSTANCE CONTAINER. manages an ISC product instance",
	Hidden: true,
	Run:    internalStart,
}

func init() {
	log.SetOutput(ioutil.Discard) // This is to silence the logging from go-plugin
	internalStartFlags.PluginFlags = make(map[string]*iscenv.PluginFlags)
	if err := addStarterFlags(internalStartCmd, &internalStartFlags.Plugins, internalStartFlags.PluginFlags); err != nil {
		app.Fatalf("%s\n", err)
	}

	internalStartCmd.Flags().StringVarP(&internalStartFlags.Instance, "instance", "i", "docker", "The instance to manage")
	internalStartCmd.Flags().StringVarP(&internalStartFlags.CControlPath, "ccontrolpath", "c", "ccontrol", "The path to the ccontrol executable in the image")

	rootCmd.AddCommand(internalStartCmd)
}

func internalStart(_ *cobra.Command, _ []string) {
	app.EnsureWithinContainer("_start")

	manager, err := app.NewInternalInstanceManager(internalStartFlags.Instance, internalStartFlags.CControlPath)
	if err != nil {
		app.Fatalf("Error creating instance manager, error: %s\n", err)
	}

	// We can't use the closing activator because we need the plugins to keep running the whole time that _start runs
	pluginsToActivate := strings.Split(internalStartFlags.Plugins, ",")
	starters := make([]iscenv.Starter, len(pluginsToActivate))
	i := 0
	pm, err := activateStarters(pluginsToActivate, func(_, _ string, starter iscenv.Starter) error {
		starters[i] = starter
		i++
		return nil
	})

	if pm != nil {
		defer pm.Close()
	}

	if err != nil {
		app.Fatalf("Failed to activate plugins, error: %s", err)
	}

	for i, starter := range starters {
		plugin := pluginsToActivate[i]
		fmt.Printf("Performing BeforeInstance step for %s\n", plugin)
		if err := starter.BeforeInstance(*manager.InternalInstanceState); err != nil {
			app.Fatalf("Failed to execute before instance hook, plugin: %s, error: %s\n", plugin, err)
		}
	}

	manager.InstanceRunningHandler = func(*iscenv.InternalInstanceState) {
		for i, starter := range starters {
			plugin := pluginsToActivate[i]
			fmt.Printf("Performing WithInstance step for %s\n", plugin)
			if err := starter.WithInstance(*manager.InternalInstanceState); err != nil {
				app.Fatalf("Failed to execute with instance hook, plugin: %s, error: %s\n", plugin, err)
			}
		}
	}
	err = manager.Manage()

	for i, starter := range starters {
		plugin := pluginsToActivate[i]
		fmt.Printf("Performing AfterInstance step for %s\n", plugin)
		if err := starter.AfterInstance(*manager.InternalInstanceState); err != nil {
			app.Fatalf("Failed to execute after instance hook, plugin: %s, error: %s\n", plugin, err)
		}
	}

	if err != nil {
		app.Fatalf("Error managing instance, error: %s\n", err)
	}
}
