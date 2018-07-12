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

package plugins

// PluginArgs holds information about the arguments for a plugin
type PluginArgs struct {
	LogLevel string
	LogJSON  bool
}

// ToArgs returns a slice of strings representing the command arguments for a plugin
func (pa PluginArgs) ToArgs() []string {
	args := []string{}
	if pa.LogLevel != "" {
		args = append(args, "--log-level="+pa.LogLevel)

		if pa.LogJSON {
			args = append(args, "--log-json=true")
		}
	}

	return args
}
