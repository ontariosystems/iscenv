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

package iscenv

import (
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

const (
	StarterKey = "start"
)

// A plugin that is executed during instance starts
type Starter interface {
	// Host hooks

	// Runs on host - Returns an array of additional flags to add to the start command.  These flags will be passed to the remaining *external* plugin hooks.  Plugin hooks within the container and expected to depend upon environment variables or volumes configured  by the host hooks.
	Flags() (PluginFlags, error)

	// Returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
	//	Environment(flags PluginFlags) ([]string, error)

	// Returns an array of volumes to add where the string is a standard docker volume format "src:dest:
	//Volumes(flags PluginFlags) ([]string, error)

	//
	//	// Additional ports to map
	//	Ports() (map[int]int, error)
	//
	//	// Will run within the container before the instance starts
	//	BeforeInstance() error
	//
	//	// Will run within the container after the instance starts
	//	WithInstance() error

	// TODO: Implement these one at a time, This belongs in a different type of plugin
	//	// Finds a list of versions and returns a map of version to full image name
	//	FindVersions() (map[string]string, error)
}

// The client (primary executable) RPC-based implementation of the interface
type StarterRPC struct{ client *rpc.Client }

func (s StarterRPC) Flags() (PluginFlags, error) {
	var resp PluginFlags
	err := s.client.Call("Plugin.Flags", new(interface{}), &resp)
	log.Println(resp)
	return resp, err
}

//func (s StarterRPC) Environment(flags PluginFlags) (map[string]string, error)

// The server (plugin side) RPC wrapper around the concrete plugin implementation
type StarterRPCServer struct{ Plugin Starter }

func (s *StarterRPCServer) Flags(args interface{}, resp *PluginFlags) (err error) {
	*resp, err = s.Plugin.Flags()
	return err
}

// The actual plugin interface needed by go-plugin.  It's a little strange in that it has both the client and server sides in the same interface.
type StarterPlugin struct {
	// The actual implementation of the plugin.  This will be unset on the client side
	Plugin Starter
}

func (s StarterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &StarterRPCServer{Plugin: s.Plugin}, nil
}

func (StarterPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &StarterRPC{client: c}, nil
}
