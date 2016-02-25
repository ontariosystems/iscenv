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

	"github.com/hashicorp/go-plugin"
)

const (
	VersionerKey = "versions"
)

// The versioner interface describes a plugin which can find versions for iscenv
type Versioner interface {
	// Find the versions avaiable for the provided image
	Versions(image string) (ISCVersions, error)
}

// The client (primary executable) RPC-based implementation of the interface
type VersionerRPC struct{ client *rpc.Client }

func (v VersionerRPC) Versions(image string) (ISCVersions, error) {
	var resp ISCVersions
	err := v.client.Call("Plugin.Versions", image, &resp)
	return resp, err
}

// The server (plugin side) RPC wrapper around the concrete plugin implementation
type VersionerRPCServer struct{ Plugin Versioner }

func (v *VersionerRPCServer) Versions(image string, resp *ISCVersions) (err error) {
	*resp, err = v.Plugin.Versions(image)
	return err
}

// The actual plugin interface needed by go-plugin.  It's a little strange in that it has both the client and server sides in the same interface.
type VersionerPlugin struct {
	// The actual implementation of the plugin.  This will be unset on the client side
	Plugin Versioner
}

func (v VersionerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &VersionerRPCServer{Plugin: v.Plugin}, nil
}

func (VersionerPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &VersionerRPC{client: c}, nil
}
