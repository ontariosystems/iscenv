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

package homedirplugin

import (
	"fmt"
	"os/user"

	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
)

const (
	pluginKey = "homedir"
)

// Plugin represents this plugin and serves as a place to attach functions to implement the Lifecycler interface
type Plugin struct{}

// Main serves as the main entry point for the plugin
func (plugin *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(plugin)
}

// Key returns the unique identifier for the plugin
func (*Plugin) Key() string {
	return pluginKey
}

// Flags returns an array of additional flags to add to the start command
func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	return iscenv.NewPluginFlags(), nil
}

// Environment returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
func (*Plugin) Environment(_ string, _ map[string]interface{}) ([]string, error) {
	home, err := getUserHome()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("HOST_HOME=%s", home),
	}, nil
}

// Copies returns an array of items to copy to the container in the format "src:dest"
func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Volumes returns an array of volumes to add where the string is a standard docker volume format "src:dest:flag"
func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	home, err := getUserHome()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("%[1]s:%[1]s:rw", home),
	}, nil
}

// Ports returns an array of additional ports to map in the format <optional hostIP>:hostPort:containerPort.  You may also prefix the host port with a + to indicate it should be shifted by the port offset
func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return []string{}, nil
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
	return nil
}

// WithInstance will run within the container after the instance starts
func (*Plugin) WithInstance(state *isclib.Instance) error {
	return nil
}

// AfterInstance will run within the container after the instance stops
func (*Plugin) AfterInstance(state *isclib.Instance) error {
	return nil
}

func getUserHome() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}

	if current.HomeDir == "" {
		return "", fmt.Errorf("Could not determine home directory, user: %s", current.Username)
	}

	return current.HomeDir, nil
}
