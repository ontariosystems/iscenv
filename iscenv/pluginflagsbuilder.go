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

import "github.com/hashicorp/go-multierror"

// NewPluginFlagsBuilder creates a new instance of object meant to assist plugins in building flags
func NewPluginFlagsBuilder() *PluginFlagsBuilder {
	return &PluginFlagsBuilder{
		flags: NewPluginFlags(),
	}
}

type PluginFlagsBuilder struct {
	flags  PluginFlags
	result *multierror.Error
}

func (builder *PluginFlagsBuilder) AddFlag(flag string, hasConfig bool, defaultValue interface{}, usage string) {
	if err := builder.flags.AddFlag(flag, hasConfig, defaultValue, usage); err != nil {
		builder.result = multierror.Append(builder.result, err)
	}
}

func (builder *PluginFlagsBuilder) Flags() (PluginFlags, error) {
	return builder.flags, builder.result.ErrorOrNil()
}
