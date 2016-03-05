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

package main

import (
	"bytes"

	"github.com/ontariosystems/iscenv/iscenv"
)

var plog = iscenv.PLog()

func main() {
	iscenv.ServeStartPlugin(new(Plugin))
}

type Plugin struct{}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	return iscenv.NewPluginFlags(), nil
}

func (*Plugin) Environment(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) BeforeInstance(instance iscenv.InternalInstance) error {
	return nil
}

func (*Plugin) WithInstance(instance iscenv.InternalInstance) error {
	plog.Debug("Enabling ISC service bindings")

	code := `MAIN
 new
 set p("Enabled")=1
 set s=##class(Security.Services).Modify("%Service_Bindings", .p)
 if $system.Status.IsError(s) {
   do $system.Status.DisplayStatus(s)
   do $zutil(4,$job,2)
 }
 quit

`
	r := bytes.NewReader([]byte(code))
	out, err := instance.Execute("%SYS", r)
	elog := plog.WithField("output", out)
	if err != nil {
		elog.WithError(err).Error("Failed to enable %Service_Bindings")
		return err
	}

	elog.Debug("Enabled service bindings")
	return nil
}

func (*Plugin) AfterInstance(instance iscenv.InternalInstance) error {
	return nil
}
