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

package plugins

import "github.com/ontariosystems/iscenv/v3/iscenv"

// NewLifecyclerManager creates and returns a PluginManager for a LifecyclerPlugin
func NewLifecyclerManager(args PluginArgs) (*LifecyclerManager, error) {
	pm, err := NewPluginManager(iscenv.LifecyclerKey, iscenv.LifecyclerPlugin{}, args)
	if err != nil {
		return nil, err
	}

	return &LifecyclerManager{PluginManager: pm}, nil
}

// LifecyclerManager is a PluginManager for managing lifecycler plugins
type LifecyclerManager struct {
	*PluginManager
}

// ActivatedLifecycler holds information about a lifecycler plugin that has been activated
type ActivatedLifecycler struct {
	*ActivatedPlugin
	Lifecycler iscenv.Lifecycler
}

// ActivatePlugins will activate the provided list of lifecycler plugins.
func (lm *LifecyclerManager) ActivatePlugins(pluginsToActivate []string) ([]*ActivatedLifecycler, error) {
	plugins, err := lm.PluginManager.ActivatePlugins(pluginsToActivate)
	if err != nil {
		return nil, err
	}

	lifecyclers := make([]*ActivatedLifecycler, len(plugins))
	for i, plugin := range plugins {
		lifecyclers[i] = &ActivatedLifecycler{
			ActivatedPlugin: plugin,
			Lifecycler:      plugin.Plugin.(iscenv.Lifecycler),
		}
	}

	return lifecyclers, nil
}
