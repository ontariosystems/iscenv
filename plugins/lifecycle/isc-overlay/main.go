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

package iscoverlayplugin

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/isclib"
	log "github.com/Sirupsen/logrus"
	"syscall"
)

var (
	plog = log.WithField("plugin", pluginKey)
)

const (
	pluginKey = "isc-overlay"
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

func (*Plugin) Environment(_ string, _ map[string]interface{}) ([]string, error) {
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
	if err := processDatsInDirectory("/ensemble"); err != nil {
		return err
	}

	if err := processDatsInDirectory("/data/db"); err != nil {
		return err
	}

	if err := processDatsInDirectory("/routine"); err != nil {
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

func processDatsInDirectory(directory string) error {
	info, err := os.Stat(directory)
	if err != nil {
		return err
	}

	stat := info.Sys().(*syscall.Stat_t)
	rootDev := stat.Dev

	return filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		stat := info.Sys().(*syscall.Stat_t)
		// We're not on the same mount point any more so skip this directory
		if stat.Dev != rootDev {
			return filepath.SkipDir
		}

		if f.IsDir() || f.Name() != "CACHE.DAT" {
			plog.WithField("path", path).Debug("Skipping non CACHE.DAT entry")
			return nil
		}
		return touchDat(path)
	})
}

func touchDat(path string) error {
	plog.WithField("path", path).Info("Touching CACHE.DAT")
	t := time.Now()
	return os.Chtimes(path, t, t)
}