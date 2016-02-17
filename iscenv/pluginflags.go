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

package iscenv

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func NewPluginFlags() PluginFlags {
	return PluginFlags{
		Flags: make(map[string]PluginFlag),
	}
}

type PluginFlags struct {
	Flags map[string]PluginFlag
}

// Add a Plugin Flag to the list of available flags.
func (pf *PluginFlags) AddFlag(flag string, defaultValue interface{}, usage string) error {
	flag = strings.ToLower(flag)
	if _, ok := pf.Flags[flag]; ok {
		return fmt.Errorf("Flag already exists, flag: %s", flag)
	}

	pf.Flags[flag] = NewPluginFlag(flag, defaultValue, usage)
	return nil
}

// This is to be called by iscenv primary process
func (pf *PluginFlags) AddFlagsToFlagSet(prefix string, flags *pflag.FlagSet) {
	for _, flag := range pf.Flags {
		flag.AddFlagToFlagSet(prefix, flags)
	}
}
