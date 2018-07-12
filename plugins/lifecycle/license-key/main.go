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

package licensekeyplugin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/isclib"
)

var plog = log.WithField("plugin", pluginKey)

const (
	pluginKey = "license-key"
	envName   = "ISC_KEY_URL"
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
	fb.AddFlag("url", true, "", "The full URL of the ISC product key to download")
	return fb.Flags()
}

// Environment returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	url, ok := flags["url"].(string)
	if !ok || url == "" {
		return nil, nil
	}

	return []string{fmt.Sprintf("%s=%s", envName, url)}, nil
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
	url := os.Getenv(envName)
	if url == "" {
		return nil
	}

	// We're to put the license key file in the manager's directory with the same owner, group and 0644 permissions
	// We'll ensure we can get the owner/group up front so we can avoid doing a bunch of work if not
	mgrDir := filepath.Join(state.DataDirectory, "mgr")
	fi, err := os.Stat(mgrDir)
	if err != nil {
		return err
	}

	stat := fi.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	gid := stat.Gid

	keyFilename := "cache.key"
	if state.Product == isclib.Iris {
		keyFilename = "license.key"
	}
	keyPath := filepath.Join(mgrDir, keyFilename)
	plog.WithFields(log.Fields{
		"url":   url,
		"path":  keyPath,
		"owner": uid,
		"group": gid,
	}).Debug("Updating ISC product license file")

	resp, err := UnsafeGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected HTTP Response, status: %s", resp.Status)
	}

	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	if _, err = io.Copy(keyFile, resp.Body); err != nil {
		return err
	}

	if err := os.Chown(keyPath, int(uid), int(gid)); err != nil {
		return err
	}

	if err := os.Chmod(keyPath, 0644); err != nil {
		return err
	}

	plog.WithFields(log.Fields{
		"url":   url,
		"path":  keyPath,
		"owner": uid,
		"group": gid,
	}).Info("Updated ISC product license file")

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
