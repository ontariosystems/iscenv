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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
	log "github.com/sirupsen/logrus"
)

var (
	plog = log.WithField("plugin", pluginKey)
)

const (
	pluginKey = "isc-overlay"
	envName   = "ISC_DAT_DIRS"
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
	fb := iscenv.NewPluginFlagsBuilder()
	fb.AddFlag("dat-dirs", true, "/ensemble,/data/db,/routine", "A comma separated list of directories in which to look for CACHE.DATs")
	return fb.Flags()
}

// Environment returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	dirs, ok := flags["dat-dirs"].(string)
	if !ok || dirs == "" {
		return nil, nil
	}

	return []string{fmt.Sprintf("%s=%s", envName, dirs)}, nil
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
	dirs := os.Getenv(envName)
	if dirs == "" {
		plog.Debug("No directories provided for DAT searching")
		return nil
	}

	for _, dir := range strings.Split(dirs, ",") {
		plog.WithField("directory", dir).Info("Processing CACHE.DAT files for directory")
		if err := processDatsInDirectory(dir); err != nil {
			return err
		}
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

func processDatsInDirectory(directory string) error {
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			plog.WithField("directory", directory).Debug("Directory doesn't exist, skipping CACHE.DAT check")
			return nil
		}
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
