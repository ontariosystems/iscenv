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

package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
)

// Config is a map of Cfgentry keyed by the Flag value
type Config map[string]Cfgentry

// Cfgentry holds the information to create a flag and a corresponding environment variable used to providing configuration to a plugin
type Cfgentry struct {
	Flag        string
	Env         string
	Description string
	Value       string
}

// Add will add a new entry to the Config keyed by the Flag value
func (c Config) Add(ce Cfgentry) {
	c[ce.Flag] = ce
}

// Get will return the Value field from the Config
func (c Config) Get(flag string) string {
	return c[flag].Value
}

// FromFlags will update the Config with the values from a map of command flags to values
func (c Config) FromFlags(flags map[string]interface{}) error {
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
					result = multierror.Append(result, fmt.Errorf("Flag value is empty, name: %s", flag))
				}
			} else {
				result = multierror.Append(result, fmt.Errorf("Flag value was not a string, name: %s, valueType: %T", flag, iv))
			}
		} else {
			result = multierror.Append(result, fmt.Errorf("Missing flag, name: %s", flag))
		}
	}
	return result
}

// FromEnv will update the Config by reading the values from the environment variables specified by the Env field
func (c Config) FromEnv() error {
	var result error

	for flag, ce := range c {
		if value, ok := os.LookupEnv(ce.Env); ok {
			ce.Value = value
		}

		if ce.Value == "" {
			result = multierror.Append(result, fmt.Errorf("environment value is empty, name: %s", flag))
		}
	}

	return result
}

// ToEnv will set the environment variables specified by the Env fields of the Config to their corresponding Value
func (c Config) ToEnv() []string {
	envs := make([]string, len(c))
	i := 0
	for _, ce := range c {
		envs[i] = fmt.Sprintf("%s=%s", ce.Env, ce.Value)
		i++
	}

	return envs
}

// Clone will clone the Config and return the copy
func (c Config) Clone() Config {
	clone := make(Config)
	for key := range c {
		clone[key] = c[key]
	}

	return clone
}
