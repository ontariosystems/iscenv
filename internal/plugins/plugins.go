/*
Copyright 2017 Ontario Systems

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

// This package exists mainly to prevent a cycle when plugins need to use "app"
package plugins

import (
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/addhostalias"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/license-key"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/isc-overlay"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/isc-source"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/csp"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/homedir"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/service-bindings"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/shm"
	"github.com/ontariosystems/iscenv/plugins/versions/local"
	"github.com/ontariosystems/iscenv/plugins/versions/quay"
	"github.com/ontariosystems/iscenv/plugins/lifecycle/journal-lck"
)

// An API for starting internal plugins
type InternalPlugin interface {
	// Start the plugin
	Main()
	Key() string
}

// A structure containing internally packaged plugins
// The first key is the "type" of plugin (versions, lifecycle, etc.)
// The second key is the suffix of the binary after (iscenv-<type>-) if the plugin were not compiled in
// The value is the implementation of the plugin itself
type internalPluginMapping map[string]map[string]InternalPlugin

var InternalPlugins internalPluginMapping

func init() {
	InternalPlugins = make(internalPluginMapping)

	addPlugin(iscenv.VersionerKey, new(localversionsplugin.Plugin))
	addPlugin(iscenv.VersionerKey, new(quayversionsplugin.Plugin))

	addPlugin(iscenv.LifecyclerKey, new(addhostaliasplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(licensekeyplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(iscoverlayplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(iscsourceplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(homedirplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(journallckplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(cspplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(servicebindingsplugin.Plugin))
	addPlugin(iscenv.LifecyclerKey, new(shmplugin.Plugin))
}

func addPlugin(pluginType string, plugin InternalPlugin) {
	if _, ok := InternalPlugins[pluginType]; !ok {
		InternalPlugins[pluginType] = make(map[string]InternalPlugin)
	}
	InternalPlugins[pluginType][plugin.Key()] = plugin
}
