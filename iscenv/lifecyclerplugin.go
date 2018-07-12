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

// Constants for lifecycler plugins
const (
	LifecyclerKey = "lifecycle"
)

// Lifecycler is an interface for plugins that is executed during instance starts
type Lifecycler interface {
	// Host hooks

	// Runs on host - Returns an array of additional flags to add to the start command.  These flags will be passed to the remaining *external* plugin hooks.  Plugin hooks within the container are expected to depend upon environment variables or volumes configured by the host hooks.
	Flags() (PluginFlags, error)

	// Returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
	Environment(version string, flagValues map[string]interface{}) ([]string, error)

	// Returns an array of items to copy to the container in the format "src:dest"
	Copies(version string, flagValues map[string]interface{}) ([]string, error)

	// Returns an array of volumes to add where the string is a standard docker volume format "src:dest:flag"
	Volumes(version string, flagValues map[string]interface{}) ([]string, error)

	// Additional ports to map in the format <optional hostIP>:hostPort:containerPort.  You may also prefix the host port with a + to indicate it should be shifted by the port offset
	Ports(version string, flagValues map[string]interface{}) ([]string, error)

	// Will run on the host after the container instance starts, receives the same flag values as start
	AfterStart(instance *ISCInstance) error

	// Will run within the container before the instance successfully starts
	BeforeInstance(state *isclib.Instance) error

	// Will run within the container after the instance starts
	WithInstance(state *isclib.Instance) error

	// Will run within the container after the instance stops
	AfterInstance(state *isclib.Instance) error

	// Will run on the host after the instance stops
	AfterStop(instance *ISCInstance) error

	// Will run on the host before the instance is removed
	BeforeRemove(instance *ISCInstance) error
}

// LifecyclerRPC is the client (primary executable) RPC-based implementation of the interface
type LifecyclerRPC struct{ client *rpc.Client }

// Flags returns an array of additional flags to add to the start command
// The logger is intentionally not passed to this method as logging cannot yet be configured during the flag setup...
func (s LifecyclerRPC) Flags() (PluginFlags, error) {
	var resp PluginFlags
	err := s.client.Call("Plugin.Flags", new(interface{}), &resp)
	return resp, err
}

// HostOpts represents information that is needed to be passed by the client to a plugin for container management
type HostOpts struct {
	Version    string
	FlagValues map[string]interface{}
}

// Environment is the client side of the plugin interface
func (s LifecyclerRPC) Environment(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Environment", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

// Copies is the client side of the plugin interface
func (s LifecyclerRPC) Copies(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Copies", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

// Volumes is the client side of the plugin interface
func (s LifecyclerRPC) Volumes(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Volumes", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

// Ports is the client side of the plugin interface
func (s LifecyclerRPC) Ports(version string, flagValues map[string]interface{}) ([]string, error) {
	var resp []string
	err := s.client.Call("Plugin.Ports", HostOpts{Version: version, FlagValues: flagValues}, &resp)
	return resp, err
}

// AfterStart is the client side of the plugin interface
func (s LifecyclerRPC) AfterStart(instance *ISCInstance) error {
	var resp struct{}
	err := s.client.Call("Plugin.AfterStart", instance, &resp)
	return err
}

// AfterStop is the client side of the plugin interface
func (s LifecyclerRPC) AfterStop(instance *ISCInstance) error {
	var resp struct{}
	err := s.client.Call("Plugin.AfterStop", instance, &resp)
	return err
}

// BeforeRemove is the client side of the plugin interface
func (s LifecyclerRPC) BeforeRemove(instance *ISCInstance) error {
	var resp struct{}
	err := s.client.Call("Plugin.BeforeRemove", instance, &resp)
	return err
}

// BeforeInstance is the client side of the plugin interface
func (s LifecyclerRPC) BeforeInstance(state *isclib.Instance) error {
	var resp struct{}
	err := s.client.Call("Plugin.BeforeInstance", state, &resp)
	return err
}

// WithInstance is the client side of the plugin interface
func (s LifecyclerRPC) WithInstance(state *isclib.Instance) error {
	var resp struct{}
	err := s.client.Call("Plugin.WithInstance", state, &resp)
	return err
}

// AfterInstance is the client side of the plugin interface
func (s LifecyclerRPC) AfterInstance(state *isclib.Instance) error {
	var resp struct{}
	err := s.client.Call("Plugin.AfterInstance", state, &resp)
	return err
}

// LifecyclerRPCServer is the server (plugin side) RPC wrapper around the concrete plugin implementation
type LifecyclerRPCServer struct{ Plugin Lifecycler }

// Flags is the server side of the plugin interface
func (s *LifecyclerRPCServer) Flags(args interface{}, resp *PluginFlags) (err error) {
	*resp, err = s.Plugin.Flags()
	return err
}

// Environment is the server side of the plugin interface
func (s *LifecyclerRPCServer) Environment(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Environment(opts.Version, opts.FlagValues)
	return err
}

// Copies is the server side of the plugin interface
func (s *LifecyclerRPCServer) Copies(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Copies(opts.Version, opts.FlagValues)
	return err
}

// Volumes is the server side of the plugin interface
func (s *LifecyclerRPCServer) Volumes(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Volumes(opts.Version, opts.FlagValues)
	return err
}

// Ports is the server side of the plugin interface
func (s *LifecyclerRPCServer) Ports(opts HostOpts, resp *[]string) (err error) {
	*resp, err = s.Plugin.Ports(opts.Version, opts.FlagValues)
	return err
}

// AfterStart is the server side of the plugin interface
func (s *LifecyclerRPCServer) AfterStart(instance *ISCInstance, resp *struct{}) (err error) {
	err = s.Plugin.AfterStart(instance)
	return err
}

// AfterStop is the server side of the plugin interface
func (s *LifecyclerRPCServer) AfterStop(instance *ISCInstance, resp *struct{}) (err error) {
	err = s.Plugin.AfterStop(instance)
	return err
}

// BeforeRemove is the server side of the plugin interface
func (s *LifecyclerRPCServer) BeforeRemove(instance *ISCInstance, resp *struct{}) (err error) {
	err = s.Plugin.BeforeRemove(instance)
	return err
}

// BeforeInstance is the server side of the plugin interface
func (s *LifecyclerRPCServer) BeforeInstance(state *isclib.Instance, resp *struct{}) (err error) {
	return s.Plugin.BeforeInstance(state)
}

// WithInstance is the server side of the plugin interface
func (s *LifecyclerRPCServer) WithInstance(state *isclib.Instance, resp *struct{}) (err error) {
	return s.Plugin.WithInstance(state)
}

// AfterInstance is the server side of the plugin interface
func (s *LifecyclerRPCServer) AfterInstance(state *isclib.Instance, resp *struct{}) (err error) {
	return s.Plugin.AfterInstance(state)
}

// LifecyclerPlugin is the actual plugin interface needed by go-plugin.  It's a little strange in that it has both the client and server sides in the same interface.
type LifecyclerPlugin struct {
	// The actual implementation of the plugin.  This will be unset on the client side
	Plugin Lifecycler
}

// Server is the server side of the plugin RPC
func (s LifecyclerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &LifecyclerRPCServer{Plugin: s.Plugin}, nil
}

// Client is the client side of the plugin RPC
func (LifecyclerPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &LifecyclerRPC{client: c}, nil
}
