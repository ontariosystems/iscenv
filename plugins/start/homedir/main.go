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
	"fmt"
	"os/user"

	"github.com/ontariosystems/iscenv/iscenv"
)

func main() {
	iscenv.ServeStartPlugin(new(Plugin))
}

type Plugin struct{}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	return iscenv.PluginFlags{}, nil
}

func (*Plugin) Environment(_ string, _ iscenv.PluginFlags) ([]string, error) {
	home, err := getUserHome()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("HOST_HOME=%s", home),
	}, nil
}

func (*Plugin) Volumes(_ string, _ iscenv.PluginFlags) ([]string, error) {
	home, err := getUserHome()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("%[1]s:%[1]s:rw", home),
	}, nil
}

func (*Plugin) Ports(_ string, _ iscenv.PluginFlags) ([]string, error) {
	return []string{}, nil
}

func (*Plugin) BeforeInstance(state iscenv.InternalInstanceState) error {
	return nil
}

func (*Plugin) WithInstance(state iscenv.InternalInstanceState) error {
	return nil
}

func (*Plugin) AfterInstance(state iscenv.InternalInstanceState) error {
	return nil
}

func getUserHome() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}

	if current.HomeDir == "" {
		return "", fmt.Errorf("Could not determine home directory, user: %s", current.Username)
	}

	return current.HomeDir, nil
}
