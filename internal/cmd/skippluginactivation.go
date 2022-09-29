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

package cmd

import (
	"os"

	"github.com/ontariosystems/iscenv/v3/iscenv"
)

func skipPluginActivation() bool {
	// If we activate during wrapped commands, it corrupts the output of wrapped commands (and plugins can fail)
	if _, wrapped := iscenv.CalledAs(); wrapped {
		return true
	}

	// If we activate during plugin calls, it infinite loops
	return len(os.Args) > 1 && os.Args[1] == pluginCmd.Use
}
