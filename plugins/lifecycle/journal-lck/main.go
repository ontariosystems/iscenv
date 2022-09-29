/*
Copyright 2017 Ontario Systems

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

package journallckplugin

import (
	"os"
	"path"

	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
	log "github.com/sirupsen/logrus"
)

var (
	plog = log.WithField("plugin", pluginKey)
)

const (
	pluginKey = "journal-lck"
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
func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Copies returns an array of items to copy to the container in the format "src:dest"
func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Volumes returns an array of volumes to add where the string is a standard docker volume format "src:dest:flag"
func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Ports returns an array of additional ports to map in the format <optional hostIP>:hostPort:containerPort.  You may also prefix the host port with a + to indicate it should be shifted by the port offset
func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
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
	plog.Info("Cleaning up primary journal directory lck files")
	pjd, err := state.DeterminePrimaryJournalDirectory()
	if err != nil {
		plog.WithError(err).Error("Failed to determine primary journal directory")
		return err
	}
	if err := deleteJournalDirectoryLck(path.Join(pjd, "cache.lck")); err != nil {
		plog.WithField("directory", pjd).WithError(err).Error("Failed to remove lck file from primary journal directory")
		return err
	}

	plog.Info("Cleaning up secondary journal directory lck files")
	sjd, err := state.DetermineSecondaryJournalDirectory()
	if err != nil {
		plog.WithError(err).Error("Failed to determine secondary journal directory")
		return err
	}
	if err := deleteJournalDirectoryLck(path.Join(sjd, "cache.lck")); err != nil {
		plog.WithField("directory", sjd).WithError(err).Error("Failed to remove lck file from secondary journal directory")
		return err
	}

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

func deleteJournalDirectoryLck(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			plog.WithField("path", path).Debug("cache.lck doesn't exist")
			return nil
		}
		return err
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	plog.WithField("path", path).Debug("cache.lck deleted")
	return nil
}
