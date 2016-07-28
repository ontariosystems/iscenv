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
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"
)

type activateLifecyclerFn func(id, pluginPath string, lifecycler iscenv.Lifecycler) error

// if pluginsToActivate is nil (rather than an empty slice, it means all)
func activateLifecyclersAndClose(pluginsToActivate []string, pluginArgs app.PluginArgs, fn activateLifecyclerFn) error {
	pm, err := activateLifecyclers(pluginsToActivate, pluginArgs, fn)
	if pm != nil {
		defer pm.Close()
	}
	return err
}

// if pluginsToActivate is nil (rather than an empty slice, it means all)
// You must check to see if pm is nil and call close *even* if there is an error
func activateLifecyclers(pluginsToActivate []string, pluginArgs app.PluginArgs, fn activateLifecyclerFn) (pm *app.PluginManager, err error) {
	pm, err = app.NewPluginManager(iscenv.ApplicationName, iscenv.LifecyclerKey, iscenv.LifecyclerPlugin{}, pluginArgs)
	if err != nil {
		return nil, err
	}

	if pluginsToActivate == nil {
		pluginsToActivate = pm.AvailablePlugins()
	}

	err = pm.ActivatePlugins(pluginsToActivate, func(id, pluginPath string, raw interface{}) error {
		lifecycler := raw.(iscenv.Lifecycler)
		return fn(id, pluginPath, lifecycler)
	})

	return pm, err
}
