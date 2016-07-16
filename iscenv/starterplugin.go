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
	"net/rpc"

	"github.com/ontariosystems/isclib"

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
	Environment(version string, flagValues map[string]interface{}) ([]string, error)

	// Returns an array of items to copy to the container in the format "src:dest"
	Copies(version string, flagValues map[string]interface{}) ([]string, error)

	// Returns an array of volumes to add where the string is a standard docker volume format "src:dest:flag"
	Volumes(version string, flagValues map[string]interface{}) ([]string, error)

	// Additional ports to map in the format <optional hostIP>:hostPort:containerPort.  You may also prefix the host port with a + to indicate it should be shifted by the port offset
	Ports(version string, flagValues map[string]interface{}) ([]string, error)

	// Will run within the container before the instance starts
	BeforeInstance(state *isclib.Instance) error

	// Will run within the container after the instance starts
	WithInstance(state *isclib.Instance) error

	// Will run within the container after the instance stops
	AfterInstance(state *isclib.Instance) error
}

// The client (primary executable) RPC-based implementation of the interface
type StarterRPC struct{ client *rpc.Client }

// The logger is intentionally not passed to this method as logging cannot yet be configured during the flag setup...
func (s StarterRPC) Flags() (PluginFlags, error) {
	var resp PluginFlags
	err := s.client.Call("Plugin.Flags", new(interface{}), &resp)
	return resp, err
}

type HostOpts struct {
	Version    string
	FlagValues map[string]interface{}
}

func (s StarterRPC) Environment(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Environment", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

func (s StarterRPC) Copies(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Copies", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

func (s StarterRPC) Volumes(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Volumes", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

func (s StarterRPC) Ports(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Ports", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

func (s StarterRPC) BeforeInstance(state *isclib.Instance) error {
	var resp struct{}
	err := s.client.Call("Plugin.BeforeInstance", state, &resp)
	return err
}

func (s StarterRPC) WithInstance(state *isclib.Instance) error {
	var resp struct{}
	err := s.client.Call("Plugin.WithInstance", state, &resp)
	return err
}

func (s StarterRPC) AfterInstance(state *isclib.Instance) error {
	var resp struct{}
	err := s.client.Call("Plugin.AfterInstance", state, &resp)
	return err
}

// The server (plugin side) RPC wrapper around the concrete plugin implementation
type StarterRPCServer struct{ Plugin Starter }

func (s *StarterRPCServer) Flags(args interface{}, resp *PluginFlags) (err error) {
	*resp, err = s.Plugin.Flags()
	return err
}

func (s *StarterRPCServer) Environment(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Environment(opts.Version, opts.FlagValues)
	return err
}

func (s *StarterRPCServer) Copies(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Volumes(opts.Version, opts.FlagValues)
	return err
}

func (s *StarterRPCServer) Volumes(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Volumes(opts.Version, opts.FlagValues)
	return err
}

func (s *StarterRPCServer) Ports(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Ports(opts.Version, opts.FlagValues)
	return err
}

func (s *StarterRPCServer) BeforeInstance(state *isclib.Instance, resp *struct{}) (err error) {
	return s.Plugin.BeforeInstance(state)
}

func (s *StarterRPCServer) WithInstance(state *isclib.Instance, resp *struct{}) (err error) {
	return s.Plugin.WithInstance(state)
}

func (s *StarterRPCServer) AfterInstance(state *isclib.Instance, resp *struct{}) (err error) {
	return s.Plugin.AfterInstance(state)
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
