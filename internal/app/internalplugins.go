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

import (
	"github.com/ontariosystems/iscenv/plugins/start/license-key"
	"github.com/ontariosystems/iscenv/plugins/start/homedir"
	"github.com/ontariosystems/iscenv/plugins/start/service-bindings"
	"github.com/ontariosystems/iscenv/plugins/start/shm"
	"github.com/ontariosystems/iscenv/plugins/versions/local"
	"github.com/ontariosystems/iscenv/plugins/versions/quay"
)

// An API for starting internal plugins
type InternalPlugin interface {
	// Start the plugin
	Main()
	Key() string
}

// A structure containing internally packaged plugins
// The first key is the "type" of plugin (versions, start, etc.)
// The second key is the suffix of the binary after (iscenv-<type>-) if the plugin were not compiled in
// The value is the implementation of the plugin itself
type internalPluginMapping map[string]map[string]InternalPlugin

var InternalPlugins internalPluginMapping

func init() {
	InternalPlugins = make(internalPluginMapping)

	addPlugin("versions", new(localversionsplugin.Plugin))
	addPlugin("versions", new(quayversionsplugin.Plugin))

	addPlugin("start", new(licensekeyplugin.Plugin))
	addPlugin("start", new(homedirplugin.Plugin))
	addPlugin("start", new(servicebindingsplugin.Plugin))
	addPlugin("start", new(shmplugin.Plugin))
}

func addPlugin(pluginType string, plugin InternalPlugin) {
	if _, ok := InternalPlugins[pluginType]; !ok {
		InternalPlugins[pluginType] = make(map[string]InternalPlugin)
	}
	InternalPlugins[pluginType][plugin.Key()] = plugin
}
