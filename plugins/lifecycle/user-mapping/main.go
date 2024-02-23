/*
Copyright 2024 Finvi, Ontario Systems

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

package usermappingplugin

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
	log "github.com/sirupsen/logrus"
)

const (
	pluginKey   = "user-mapping"
	mappingFlag = "users"
	ownerFlag   = "owner"
	envMapping  = "ISCENV_USER_MAPPING"
	envOwner    = "ISCENV_REMAP_OWNER"
)

var plog = log.WithField("plugin", pluginKey)

type Plugin struct{}

// Main serves as the main entry point for the plugin
func (p *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(p)
}

// Key returns the unique identifier for the plugin
func (*Plugin) Key() string {
	return pluginKey
}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	fb := iscenv.NewPluginFlagsBuilder()
	fb.AddFlag(mappingFlag, true, "", "Comma separated list of username:group:uid:gid to replace in instance")
	fb.AddFlag(ownerFlag, true, true, "Should the ISC instance owner user be remapped to the current user")

	return fb.Flags()
}

func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	mapping, ok := flags[mappingFlag].(string)
	if !ok {
		return nil, fmt.Errorf("%s is not a string", mappingFlag)
	}
	env := []string{fmt.Sprintf("%s=%s", envMapping, mapping)}

	remapOwner, ok := flags[ownerFlag].(bool)
	if !ok {
		return nil, fmt.Errorf("%s is not a bool", ownerFlag)
	}

	if remapOwner {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		env = append(env, fmt.Sprintf("%s=%s:%s", envOwner, u.Uid, u.Gid))
	}

	return env, nil
}

func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) { return nil, nil }

func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) AfterStart(_ *iscenv.ISCInstance) error {
	return nil
}

func (*Plugin) AfterStop(_ *iscenv.ISCInstance) error {
	return nil
}

func (*Plugin) BeforeRemove(_ *iscenv.ISCInstance) error {
	return nil
}

func (p *Plugin) BeforeInstance(state *isclib.Instance) error {
	mappingEnv := os.Getenv(envMapping)
	replacements := []userInfo{}
	if mappingEnv != "" {
		mappings := strings.Split(mappingEnv, ",")
		for _, m := range mappings {
			parts := strings.Split(m, ":")
			if len(parts) != 4 {
				return fmt.Errorf("user mapping has wrong number of parts. mapping: '%s'", m)
			}

			uid, err := strconv.Atoi(parts[2])
			if err != nil {
				return err
			}

			gid, err := strconv.Atoi(parts[3])
			if err != nil {
				return err
			}

			replacements = append(replacements, userInfo{
				user:  parts[0],
				group: parts[1],
				uid:   uid,
				gid:   gid,
			})
		}
	}

	remapOwnerEnv := os.Getenv(envOwner)
	if remapOwnerEnv != "" {
		parts := strings.Split(remapOwnerEnv, ":")
		if len(parts) != 2 {
			return fmt.Errorf("%s has invalid value", envOwner)
		}

		hostUid, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		hostGid, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		ownerUser, ownerGroup, err := state.DetermineOwner()
		if err != nil {
			return err
		}

		replacements = append(replacements, userInfo{
			user:  ownerUser,
			group: ownerGroup,
			uid:   hostUid,
			gid:   hostGid,
		})
	}

	for _, ui := range replacements {
		if err := replaceUser(ui); err != nil {
			return err
		}
	}

	return nil
}

func (*Plugin) WithInstance(_ *isclib.Instance) error {
	return nil
}

func (*Plugin) AfterInstance(_ *isclib.Instance) error {
	return nil
}
