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

package app

import "github.com/ontariosystems/iscenv/iscenv"

func NewVersionerManager(args PluginArgs) (*VersionerManager, error) {
	pm, err := NewPluginManager(iscenv.VersionerKey, iscenv.VersionerPlugin{}, args)
	if err != nil {
		return nil, err
	}

	return &VersionerManager{PluginManager: pm}, nil
}

type VersionerManager struct {
	*PluginManager
}

type ActivatedVersioner struct {
	*ActivatedPlugin
	Versioner iscenv.Versioner
}

func (lm *VersionerManager) ActivatePlugins(pluginsToActivate []string) ([]*ActivatedVersioner, error) {
	plugins, err := lm.PluginManager.ActivatePlugins(pluginsToActivate)
	if err != nil {
		return nil, err
	}

	versioners := make([]*ActivatedVersioner, len(plugins))
	for i, plugin := range plugins {
		versioners[i] = &ActivatedVersioner{
			ActivatedPlugin: plugin,
			Versioner:       plugin.Plugin.(iscenv.Versioner),
		}
	}

	return versioners, nil
}
