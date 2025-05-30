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

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/kardianos/osext"
	"github.com/ontariosystems/iscenv/v3/iscenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Throw away the logs from go-plugin
	hclog.DefaultOutput = io.Discard
}

// NewPluginManager creates and returns a PluginManager for the requested Plugin
func NewPluginManager(pluginType string, iscenvPlugin plugin.Plugin, args PluginArgs) (*PluginManager, error) {
	exeDir, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}

	thisExe, err := osext.Executable()
	if err != nil {
		return nil, err
	}

	// The internal plugins are the defaults
	plugins := make(map[string]string)
	if internal, ok := InternalPlugins[pluginType]; ok {
		for key := range internal {
			plugins[key] = ""
		}
	}
	log.Debugf("Found %d internal %s plugin(s)", len(InternalPlugins[pluginType]), pluginType)

	log.Debugf("Searching %s for external %s plugins", pluginType, exeDir)
	exes, err := filepath.Glob(filepath.Join(exeDir, fmt.Sprintf("%s-%s-*", iscenv.ApplicationName, pluginType)))
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %d external %s plugin(s)", len(exes), pluginType)

	// Plugins on disk override the internal plugins
	for _, exe := range exes {
		key := strings.SplitN(filepath.Base(exe), "-", 3)[2]
		plugins[key] = exe
	}

	clients := make(map[string]*PluginClient)
	for key, exe := range plugins {
		var cmd *exec.Cmd
		if exe != "" {
			cmd = exec.Command(exe, args.ToArgs()...)
		} else {
			cmd = exec.Command(thisExe, append([]string{"plugin", pluginType, key}, args.ToArgs()...)...)
		}

		client := &PluginClient{
			ExecutablePath: exe,
			Client: plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: iscenv.PluginHandshake,
				Plugins:         map[string]plugin.Plugin{pluginType: iscenvPlugin},
				Cmd:             cmd,
				SyncStdout:      os.Stdout,
				SyncStderr:      os.Stderr,
			}),
		}

		clients[key] = client
	}
	log.Debugf("Found %d unique %s plugin(s)", len(clients), pluginType)

	return &PluginManager{
		pluginType: pluginType,
		clients:    clients,
	}, nil
}

// PluginManager holds the type of the plugin and a map of the clients
type PluginManager struct {
	pluginType string
	clients    map[string]*PluginClient
}

// PluginClient holds path to the executable and client for a plugin
type PluginClient struct {
	ExecutablePath string
	*plugin.Client
}

// ActivatedPlugin holds information about a plugin that has been activated
type ActivatedPlugin struct {
	Id     string
	Path   string
	Plugin interface{}
}

// RPCClient is needed because the embedded struct is Client, and it has a function called Client so it's client.Client() is ambiguous
func (pc *PluginClient) RPCClient(ctx context.Context) (plugin.ClientProtocol, error) {
	return pc.Client.Client(ctx)
}

// AvailablePlugins returns a slice of the keys of all the discovered plugins
func (pm *PluginManager) AvailablePlugins() []string {
	plugins := make([]string, len(pm.clients))
	i := 0
	for plugin := range pm.clients {
		plugins[i] = plugin
		i++
	}

	return plugins
}

// ActivatePlugins will activate the provided list of plugins.  If the list is nil, it will activate all the plugins.
// It does this by traversing all the plugins dispensing them to the rpc client and then returning an object containing the id of the plugin, the path to the executable (if not internal) and the raw plugin interface{} which the caller will likely want to typecast into something more useful.
// It will return the ActivatedPlugins in the same order as the pluginsToActivate and any error encountered
func (pm *PluginManager) ActivatePlugins(ctx context.Context, pluginsToActivate []string) ([]*ActivatedPlugin, error) {
	if pluginsToActivate == nil {
		pluginsToActivate = pm.AvailablePlugins()
	}

	activatedPlugins := make([]*ActivatedPlugin, len(pluginsToActivate))
	for i, key := range pluginsToActivate {
		key = strings.ToLower(key)

		client, ok := pm.clients[key]
		if !ok {
			return nil, fmt.Errorf("no such plugin, name: %s", key)
		}

		rpcClient, err := client.RPCClient(ctx)
		if err != nil {
			return nil, err
		}

		raw, err := rpcClient.Dispense(pm.pluginType)
		if err != nil {
			return nil, err
		}

		activatedPlugins[i] = &ActivatedPlugin{
			Id:     key,
			Path:   client.ExecutablePath,
			Plugin: raw,
		}
	}

	return activatedPlugins, nil
}

// Close will shut down all the PluginClients
func (pm *PluginManager) Close(ctx context.Context) {
	for _, client := range pm.clients {
		client.Kill(ctx)
	}
}
