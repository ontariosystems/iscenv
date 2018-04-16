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

	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/isclib"
)

var (
	plog = log.WithField("plugin", pluginKey)
)

const (
	pluginKey = "journal-lck"
)

type Plugin struct{}

func (plugin *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(plugin)
}

func (*Plugin) Key() string {
	return pluginKey
}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	return iscenv.NewPluginFlags(), nil
}

func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) AfterStart(instance *iscenv.ISCInstance) error {
	return nil
}

func (*Plugin) AfterStop(instance *iscenv.ISCInstance) error {
	return nil
}

func (*Plugin) BeforeRemove(instance *iscenv.ISCInstance) error {
	return nil
}

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

func (*Plugin) WithInstance(state *isclib.Instance) error {
	return nil
}

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
