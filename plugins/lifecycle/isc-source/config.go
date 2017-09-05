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
	"fmt"
	"os"

	multierror "github.com/hashicorp/go-multierror"
)

type config map[string]cfgentry
type cfgentry struct {
	Flag        string
	Env         string
	Description string
	Value       string
}

func (c config) Add(ce cfgentry) {
	c[ce.Flag] = ce
}

func (c config) Get(flag string) string {
	return c[flag].Value
}

func (c config) FromFlags(flags map[string]interface{}) error {
	var result error

	for flag, ce := range c {
		iv, ok := flags[flag]
		if ok {
			value, ok := iv.(string)
			if ok {
				if value != "" {
					ce.Value = value
					c[flag] = ce
				} else {
					multierror.Append(result, fmt.Errorf("Flag value is empty, name: %s", flag))
				}
			} else {
				multierror.Append(result, fmt.Errorf("Flag value was not a string, name: %s, valueType: %T", flag, iv))
			}
		} else {
			multierror.Append(result, fmt.Errorf("Missing flag, name: %s", flag))
		}
	}
	return result
}

func (c config) FromEnv() error {
	var result error

	for flag, ce := range c {
		value := os.Getenv(ce.Env)
		if value != "" {
			ce.Value = value
			c[flag] = ce
		} else {
			multierror.Append(result, fmt.Errorf("Environment value is empty, name: %s", flag))
		}
	}

	return result
}

func (c config) ToEnv() []string {
	envs := make([]string, len(c))
	i := 0
	for _, ce := range c {
		envs[i] = fmt.Sprintf("%s=%s", ce.Env, ce.Value)
		i++
	}

	return envs
}

func (c config) Clone() config {
	clone := make(config)
	for key := range c {
		clone[key] = c[key]
	}

	return clone
}
