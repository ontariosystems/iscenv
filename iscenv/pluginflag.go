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
	"log"
	"net"
	"reflect"
	"time"

	"github.com/spf13/pflag"
)

func NewPluginFlag(flag string, defaultValue interface{}, usage string) PluginFlag {
	return PluginFlag{
		Flag:         flag,
		Usage:        usage,
		DefaultValue: defaultValue,
	}
}

type PluginFlag struct {
	Flag         string
	Usage        string
	DefaultValue interface{}
	// The value field is a little wonky as the pointer will not pass across the gob (it always ends up being just the value)
	Value interface{}
}

func (pf *PluginFlag) AddFlagToFlagSet(prefix string, flags *pflag.FlagSet) error {
	pf.Value = getEmptyValuePointer(pf.DefaultValue)
	usage := fmt.Sprintf("%s (plugin: %s)", pf.Usage, prefix)
	// the "count" type is intentionally not supported as it's not *that* useful and it cannot be distinguished from an int
	// by simply using its default value.  it seemed much more important to maintain the simple API
	switch v := pf.DefaultValue.(type) {
	case pflag.Value:
		flags.Var(pf.Value.(pflag.Value), fullFlag(prefix, pf.Flag), usage)
	case bool:
		flags.BoolVar(pf.Value.(*bool), fullFlag(prefix, pf.Flag), v, usage)
	case time.Duration:
		flags.DurationVar(pf.Value.(*time.Duration), fullFlag(prefix, pf.Flag), v, usage)
	case float32:
		flags.Float32Var(pf.Value.(*float32), fullFlag(prefix, pf.Flag), v, usage)
	case float64:
		flags.Float64Var(pf.Value.(*float64), fullFlag(prefix, pf.Flag), v, usage)
	case int:
		flags.IntVar(pf.Value.(*int), fullFlag(prefix, pf.Flag), v, usage)
	case int8:
		flags.Int8Var(pf.Value.(*int8), fullFlag(prefix, pf.Flag), v, usage)
	case int32:
		flags.Int32Var(pf.Value.(*int32), fullFlag(prefix, pf.Flag), v, usage)
	case int64:
		flags.Int64Var(pf.Value.(*int64), fullFlag(prefix, pf.Flag), v, usage)
	case []int:
		flags.IntSliceVar(pf.Value.(*[]int), fullFlag(prefix, pf.Flag), v, usage)
	case net.IPMask:
		flags.IPMaskVar(pf.Value.(*net.IPMask), fullFlag(prefix, pf.Flag), v, usage)
	case net.IPNet:
		flags.IPNetVar(pf.Value.(*net.IPNet), fullFlag(prefix, pf.Flag), v, usage)
	case string:
		flags.StringVar(pf.Value.(*string), fullFlag(prefix, pf.Flag), v, usage)
	case []string:
		flags.StringSliceVar(pf.Value.(*[]string), fullFlag(prefix, pf.Flag), v, usage)
	case uint:
		flags.UintVar(pf.Value.(*uint), fullFlag(prefix, pf.Flag), v, usage)
	case uint8:
		flags.Uint8Var(pf.Value.(*uint8), fullFlag(prefix, pf.Flag), v, usage)
	case uint16:
		flags.Uint16Var(pf.Value.(*uint16), fullFlag(prefix, pf.Flag), v, usage)
	case uint32:
		flags.Uint32Var(pf.Value.(*uint32), fullFlag(prefix, pf.Flag), v, usage)
	case uint64:
		flags.Uint64Var(pf.Value.(*uint64), fullFlag(prefix, pf.Flag), v, usage)
	default:
		return fmt.Errorf("Cannot handle type: %T", pf.Value)
	}

	return nil
}

func fullFlag(prefix, flag string) string {
	return fmt.Sprintf("%s-%s", prefix, flag)
}

func getEmptyValuePointer(value interface{}) interface{} {
	return reflect.New(reflect.TypeOf(value)).Interface()
}
