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

// This package is a simple cobra/viper-aware flag wrapper that creates a unified interface for flag addition and access.
// It supports only the types that iscenv needs, if you need more add them.
package flags

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type flag struct {
	config bool
	value  interface{}
}

var flags map[string]*flag = make(map[string]*flag)

func GetKeys() []string {
	keys := make([]string, len(flags))
	i := 0
	for key := range flags {
		keys[i] = key
		i++
	}

	return keys
}

func GetString(cmd *cobra.Command, name string) string {
	return GetValue(cmd, name).(string)
}

func GetBool(cmd *cobra.Command, name string) bool {
	return GetValue(cmd, name).(bool)
}

func GetInt(cmd *cobra.Command, name string) int {
	return int(forceFloat64(GetValue(cmd, name)))
}

func GetInt64(cmd *cobra.Command, name string) int64 {
	return int64(forceFloat64(GetValue(cmd, name)))
}

func GetUint(cmd *cobra.Command, name string) uint {
	return uint(forceFloat64(GetValue(cmd, name)))
}

func GetStringSlice(cmd *cobra.Command, name string) []string {
	return GetValue(cmd, name).([]string)
}

func GetValue(cmd *cobra.Command, name string) interface{} {
	return GetValueWithKey(GetFlagKey(cmd, name))
}

func GetValueWithKey(key string) interface{} {
	flag, ok := flags[key]
	if !ok {
		panic("Attempt to access non-existent flag, key: " + key)
	}

	if flag.config {
		return viper.Get(key)
	}

	rv := reflect.ValueOf(flag.value)
	if rv.Kind() != reflect.Ptr {
		panic("They value stored at this key is not a pointer (should not be possible), key: " + key)
	}

	return rv.Elem().Interface()
}

func GetRawValue(key string) interface{} {
	flag, ok := flags[key]
	if !ok {
		panic("Attempt to access non-existent flag, key: " + key)
	}

	return flag.value
}

func AddFlag(cmd *cobra.Command, name string, value interface{}, usage string) {
	AddFlagComplex(cmd, false, false, name, "", value, usage)
}

func AddFlagP(cmd *cobra.Command, name string, shorthand string, value interface{}, usage string) {
	AddFlagComplex(cmd, false, false, name, shorthand, value, usage)
}

func AddConfigFlag(cmd *cobra.Command, name string, value interface{}, usage string) {
	AddFlagComplex(cmd, false, true, name, "", value, usage)
}

func AddConfigFlagP(cmd *cobra.Command, name string, shorthand string, value interface{}, usage string) {
	AddFlagComplex(cmd, false, true, name, shorthand, value, usage)
}

// Add a flag to the provided cobra command.  The default value will determine the type of flag.  Supported types are:
// string, bool, int64, uint, []string
func AddFlagComplex(cmd *cobra.Command, persistent bool, cfg bool, name string, shorthand string, value interface{}, usage string) {
	key := GetFlagKey(cmd, name)
	if cfg {
		usage += " (config path: " + key + ")"
	}

	f := &flag{config: cfg}
	flags[key] = f

	var cmdFlags *pflag.FlagSet
	if persistent {
		cmdFlags = cmd.PersistentFlags()
	} else {
		cmdFlags = cmd.Flags()
	}

	// There are many ways this could have been done, i think this ends up being the simplest with a minor performance hit
	switch v := value.(type) {
	case string:
		f.value = cmdFlags.StringP(name, shorthand, v, usage)
	case bool:
		f.value = cmdFlags.BoolP(name, shorthand, v, usage)
	case int:
		f.value = cmdFlags.IntP(name, shorthand, v, usage)
	case int64:
		f.value = cmdFlags.Int64P(name, shorthand, v, usage)
	case uint:
		f.value = cmdFlags.UintP(name, shorthand, v, usage)
	case []string:
		if cfg {
			panic("[]string is not support with a configuration option.  This is a limitation of viper, change your code")
		}
		f.value = cmdFlags.StringSliceP(name, shorthand, v, usage)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

	if cfg {
		if err := viper.BindPFlag(key, cmdFlags.Lookup(name)); err != nil {
			panic(err)
		}
	}
}

func GetFlagKey(cmd *cobra.Command, name string) string {
	if cmd.Name() != "" {
		return cmd.Name() + "." + name
	}

	return name
}

// viper uses spf13's cast library to quietly covert types in all manner of ways, (strings to ints by parsing for example).  We don't want to go that far, but since viper stores all numeric values as float64 we have to do a little conversion.
func forceFloat64(num interface{}) float64 {
	switch n := num.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int64:
		return float64(n)
	case int32:
		return float64(n)
	case int16:
		return float64(n)
	case int8:
		return float64(n)
	case int:
		return float64(n)
	case uint64:
		return float64(n)
	case uint32:
		return float64(n)
	case uint16:
		return float64(n)
	case uint8:
		return float64(n)
	case uint:
		return float64(n)
	}

	panic(fmt.Sprintf("unsupported type: %T", num))
}
