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

	"github.com/ontariosystems/iscenv/iscenv"
	plugin "github.com/ontariosystems/iscenv/internal/plugins/config"
	"github.com/ontariosystems/isclib"
	log "github.com/Sirupsen/logrus"
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

type Plugin struct{}

func (plugin *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(plugin)
}

func (*Plugin) Key() string {
	return pluginKey
}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	fb := iscenv.NewPluginFlagsBuilder()
	for _, ce := range cfg {
		fb.AddFlag(ce.Flag, true, ce.Value, ce.Description)
	}

	return fb.Flags()
}

func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	c := cfg.Clone()
	if err := c.FromFlags(flags); err != nil {
		return nil, err
	}

	return c.ToEnv(), nil
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

func (*Plugin) BeforeInstance(instance *isclib.Instance) error {
	return nil
}

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
		} else {
			return err
		}
	}

	if err := instance.ExecuteAsManager(); err != nil {
		return err
	}

	if err := importSource(l, instance, srcDir, ns, opts); err != nil {
		return err
	}

	return nil
}

func (*Plugin) AfterInstance(instance *isclib.Instance) error {
	return nil
}
