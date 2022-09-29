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
	"strings"

	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	"github.com/spf13/cobra"
)

func getPluginsToActivate(cmd *cobra.Command) []string {
	s := flags.GetString(cmd, "plugins")
	if strings.TrimSpace(s) == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}
