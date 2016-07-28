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

package shmplugin

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/isclib"
)

var plog = log.WithField("plugin", pluginKey)

const (
	pluginKey = "shm"
	envName   = "ISCENV_SHM_SIZE"
)

type Plugin struct{}

func (plugin *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(plugin)
}

func (*Plugin) Key() string {
	return pluginKey
}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	fb := iscenv.NewPluginFlagsBuilder()
	fb.AddFlag("size", true, 8192, "The size of shared memory in megabytes")
	return fb.Flags()
}

func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	shmsize, ok := flags["size"].(int)
	if !ok || shmsize == 0 {
		return nil, nil
	}

	return []string{fmt.Sprintf("%s=%d", envName, shmsize)}, nil
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

func (*Plugin) AfterStart(version string, flagValues map[string]interface{}) error {
	return nil
}

func (*Plugin) AfterStop(version string, flagValues map[string]interface{}) error {
	return nil
}

func (*Plugin) AfterRemove(version string, flagValues map[string]interface{}) error {
	return nil
}

func (*Plugin) BeforeInstance(state *isclib.Instance) error {
	sizeStr := os.Getenv(envName)
	if sizeStr == "" {
		return nil
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return err
	}

	sizeBytes := size * 1024 * 102

	if err := sysctl("shmmax", sizeBytes); err != nil {
		return err
	}

	if err := sysctl("shmall", sizeBytes); err != nil {
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

func sysctl(param string, size int) error {
	if b, err := exec.Command("sysctl", "-w", fmt.Sprintf("kernel.%s=%d", param, size)).CombinedOutput(); err != nil {
		out := string(b)
		plog.WithFields(log.Fields{
			"param":  param,
			"size":   size,
			"stdout": out,
		}).WithError(err).Error("Failed to set shared memory")
		return err
	}

	return nil
}
