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
	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/internal/iscenv"
)

var csessionFlags = struct {
	Namespace string
}{}

var csessionCmd = &cobra.Command{
	Use:   "csession INSTANCE",
	Short: "Start csession for instance",
	Long:  "Connect to an instance with SSH using private key auth and initiate a csession.",
	Run:   csession,
}

func init() {
	rootCmd.AddCommand(csessionCmd)

	csessionCmd.Flags().StringVarP(&csessionFlags.Namespace, "namespace", "n", "%SYS", "Use a specific staring namespace")
}

func csession(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		csession := "csession docker"
		if csessionFlags.Namespace != "" {
			csession += " -U" + csessionFlags.Namespace
		}
		sshExec(args[0], nil, csession)
	} else {
		iscenv.Fatal("Must provide an instance")
	}
}
