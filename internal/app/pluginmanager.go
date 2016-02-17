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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/kardianos/osext"
	"github.com/ontariosystems/iscenv/iscenv"
)

// TODO: It's likely that this could be written more cleanly

type pluginVisitorFn func(id string, raw interface{}) error

func NewPluginManager(applicationName, pluginType string, iscenvPlugin plugin.Plugin) (*PluginManager, error) {
	exeDir, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}

	exes, err := filepath.Glob(filepath.Join(exeDir, fmt.Sprintf("%s-%s-*", applicationName, pluginType)))
	if err != nil {
		return nil, err
	}

	clients := make(map[string]*plugin.Client)
	for _, exe := range exes {
		id := strings.SplitN(filepath.Base(exe), "-", 3)[2]
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: iscenv.PluginHandshake,
			Plugins:         map[string]plugin.Plugin{pluginType: iscenvPlugin},
			Cmd:             exec.Command(exe),
		})

		clients[id] = client
	}

	return &PluginManager{
		pluginType: pluginType,
		clients:    clients,
	}, nil
}

type PluginManager struct {
	pluginType string
	clients    map[string]*plugin.Client
}

// This will traverse all of the plugins dispensing them to the rpc client and then returning the raw interface{} returns, the caller will want to type cast it to the appropriate interface
func (pm *PluginManager) VisitPlugins(fn pluginVisitorFn) error {
	for key, client := range pm.clients {
		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		raw, err := rpcClient.Dispense(pm.pluginType)
		if err != nil {
			return err
		}

		err = fn(key, raw)
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
