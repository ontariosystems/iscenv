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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/ontariosystems/iscenv/iscenv"
)

var plog = iscenv.PLog()

const (
	envName = "ISC_KEY_URL"
)

func main() {
	iscenv.ServeStartPlugin(new(Plugin))
}

type Plugin struct{}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	flags := iscenv.NewPluginFlags()
	flags.AddFlag("url", true, "", "The full URL of the ISC product key to download")
	return flags, nil
}

func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	url, ok := flags["url"].(string)
	if !ok || url == "" {
		return nil, nil
	}

	return []string{fmt.Sprintf("%s=%s", envName, url)}, nil
}

func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) BeforeInstance(state iscenv.InternalInstance) error {
	url := os.Getenv(envName)
	if url == "" {
		return nil
	}

	// We're to put the license key file in the manager's directory with the same owner, group and 0644 permissions
	// We'll ensure we can get the owner/group up front so we can avoid doing a bunch of work if not
	mgrDir := filepath.Join(state.Directory, "mgr")
	fi, err := os.Stat(mgrDir)
	if err != nil {
		return err
	}

	stat := fi.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	gid := stat.Gid

	keyPath := filepath.Join(mgrDir, "cache.key")
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

	os.Chown(keyPath, int(uid), int(gid))
	os.Chmod(keyPath, 0644)

	plog.WithFields(log.Fields{
		"url":   url,
		"path":  keyPath,
		"owner": uid,
		"group": gid,
	}).Info("Updated ISC product license file")

	return nil
}

func (*Plugin) WithInstance(state iscenv.InternalInstance) error {
	return nil
}

func (*Plugin) AfterInstance(state iscenv.InternalInstance) error {
	return nil
}
