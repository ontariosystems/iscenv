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

package iscsourceplugin

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	plugin "github.com/ontariosystems/iscenv/config"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/isclib"
)

const (
	pluginKey  = "isc-source"
	srcFlag    = "src-dir"
	nsFlag     = "namespace"
	importFlag = "import-options"
)

var (
	plog = log.WithField("plugin", pluginKey)
	cfg  plugin.Config
)

func init() {
	cfg = make(plugin.Config)
	cfg.Add(plugin.Cfgentry{
		Flag:        srcFlag,
		Env:         "ISC_SRC_DIR",
		Description: "The directory containing the distribution files",
		Value:       "/isc_src",
	})
	cfg.Add(plugin.Cfgentry{
		Flag:        nsFlag,
		Env:         "CORE_NAMESPACE",
		Description: "The existing namespace into which the source will be imported",
		Value:       "USER",
	})
	cfg.Add(plugin.Cfgentry{
		Flag:        importFlag,
		Env:         "IMPORT_OPTIONS",
		Description: "The options to use for the source being imported",
		Value:       "/subclasses/compile/percent/predecessorclasses/relatedclasses",
	})
}

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
	for _, ce := range cfg {
		fb.AddFlag(ce.Flag, true, ce.Value, ce.Description)
	}

	return fb.Flags()
}

// Environment returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	c := cfg.Clone()
	if err := c.FromFlags(flags); err != nil {
		return nil, err
	}

	return c.ToEnv(), nil
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
func (*Plugin) BeforeInstance(instance *isclib.Instance) error {
	return nil
}

// WithInstance will run within the container after the instance starts
func (*Plugin) WithInstance(instance *isclib.Instance) error {
	l := plog
	l.Debug("Executing source loading process")

	c := cfg.Clone()
	if err := c.FromEnv(); err != nil {
		return err
	}

	srcDir := c.Get(srcFlag)
	ns := strings.ToUpper(c.Get(nsFlag))
	opts := c.Get(importFlag)

	l = l.WithFields(log.Fields{
		"srcDir":         srcDir,
		"namespace":      ns,
		"import-options": opts,
	})

	if _, err := os.Stat(srcDir); err != nil {
		if os.IsNotExist(err) {
			l.Warn("Distribution directory does not exist, skipping core install")
			return nil
		}

		return err
	}

	if err := instance.ExecuteAsManager(); err != nil {
		return err
	}

	return importSource(l, instance, srcDir, ns, opts)
}

// AfterInstance will run within the container after the instance stops
func (*Plugin) AfterInstance(instance *isclib.Instance) error {
	return nil
}
