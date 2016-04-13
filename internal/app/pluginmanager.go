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
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-plugin"
	"github.com/kardianos/osext"
	"github.com/ontariosystems/iscenv/iscenv"
)

type activatePluginFn func(id, pluginPath string, raw interface{}) error

func init() {
	// Throw away the logs from go-plugin
	golog.SetOutput(ioutil.Discard)
}

func NewPluginManager(applicationName, pluginType string, iscenvPlugin plugin.Plugin, args PluginArgs) (*PluginManager, error) {
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
	exes, err := filepath.Glob(filepath.Join(exeDir, fmt.Sprintf("%s-%s-*", applicationName, pluginType)))
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

type PluginManager struct {
	pluginType string
	clients    map[string]*PluginClient
}

type PluginClient struct {
	ExecutablePath string
	*plugin.Client
}

// Needed because the embedded struct is Client and it has a function called Client so it's client.Client() is ambiguous
func (pc *PluginClient) RPCClient() (*plugin.RPCClient, error) {
	return pc.Client.Client()
}

func (pm *PluginManager) AvailablePlugins() []string {
	plugins := make([]string, len(pm.clients))
	i := 0
	for plugin := range pm.clients {
		plugins[i] = plugin
		i++
	}

	return plugins
}

// This will traverse all of the plugins dispensing them to the rpc client and then returning the raw interface{} returns, the caller will want to type cast it to the appropriate interface
func (pm *PluginManager) ActivatePlugins(pluginsToActivate []string, fn activatePluginFn) error {
	for _, key := range pluginsToActivate {
		key = strings.ToLower(key)

		client, ok := pm.clients[key]
		if !ok {
			return fmt.Errorf("No such plugin, name: %s", key)
		}

		rpcClient, err := client.RPCClient()
		if err != nil {
			return err
		}

		raw, err := rpcClient.Dispense(pm.pluginType)
		if err != nil {
			return err
		}

		err = fn(key, client.ExecutablePath, raw)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pm *PluginManager) Close() {
	for _, client := range pm.clients {
		client.Kill()
	}
}
