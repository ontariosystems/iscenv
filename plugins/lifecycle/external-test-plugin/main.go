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

package main

import (
	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
	log "github.com/sirupsen/logrus"
)

var plog = log.WithField("plugin", "external-test-plugin")

// Plugin represents this plugin and serves as a place to attach functions to implement the Lifecycler interface
type Plugin struct{}

func main() {
	iscenv.ServeLifecyclePlugin(new(Plugin))
}

// Flags returns an array of additional flags to add to the start command
func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	plog.Info("Flags")
	return iscenv.NewPluginFlags(), nil
}

// Environment returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
func (*Plugin) Environment(version string, flags map[string]interface{}) ([]string, error) {
	plog.Info("Environment")
	return nil, nil
}

// Copies returns an array of items to copy to the container in the format "src:dest"
func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Volumes returns an array of volumes to add where the string is a standard docker volume format "src:dest:flag"
func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	plog.Info("Volumes")
	return nil, nil
}

// Ports returns an array of additional ports to map in the format <optional hostIP>:hostPort:containerPort.  You may also prefix the host port with a + to indicate it should be shifted by the port offset
func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	plog.Info("Ports")
	return nil, nil
}

// AfterStart will run on the host after the container instance starts, receives the same flag values as start
func (*Plugin) AfterStart(instance *iscenv.ISCInstance) error {
	return nil
}

// AfterStop will run on the host after the instance stops
func (*Plugin) AfterStop(instance *iscenv.ISCInstance) error {
	return nil
}

// BeforeRemove will run on the host before the instance is removed
func (*Plugin) BeforeRemove(instance *iscenv.ISCInstance) error {
	return nil
}

// BeforeInstance will run within the container before the instance successfully starts
func (*Plugin) BeforeInstance(state *isclib.Instance) error {
	plog.Info("BeforeInstance")
	return nil
}

// WithInstance will run within the container after the instance starts
func (*Plugin) WithInstance(state *isclib.Instance) error {
	plog.Info("WithInstance")
	return nil
}

// AfterInstance will run within the container after the instance stops
func (*Plugin) AfterInstance(state *isclib.Instance) error {
	plog.Info("AfterInstance")
	return nil
}
