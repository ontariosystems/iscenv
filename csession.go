/*
Copyright 2014 Ontario Systems

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
	"github.com/spf13/cobra"
)

var csessionNamespace string

var csessionCommand = &cobra.Command{
	Use:   "csession INSTANCE",
	Short: "Start csession for instance",
	Long:  "Connect to an instance with SSH using private key auth and initiate a csession.",
}

func init() {
	csessionCommand.Run = csession
	csessionCommand.Flags().StringVarP(&csessionNamespace, "namespace", "n", "%SYS", "Use a specific staring namespace")
}

func csession(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		csession := "csession docker"
		if csessionNamespace != "" {
			csession += " -U" + csessionNamespace
		}
		sshExec(args[0], nil, csession)
	} else {
		Fatal("Must provide an instance")
	}
}
