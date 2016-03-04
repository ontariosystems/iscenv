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

import "fmt"

func NewPluginFlag(flag string, defaultValue interface{}, usage string) *PluginFlag {
	return &PluginFlag{
		Flag:         flag,
		Usage:        usage,
		DefaultValue: defaultValue,
	}
}

type PluginFlag struct {
	Flag         string
	HasConfig    bool
	DefaultValue interface{}
	Usage        string
}

func fullFlag(prefix, flag string) string {
	return fmt.Sprintf("%s-%s", prefix, flag)
}
