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

	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func validateMultiInstanceFlags(cmd *cobra.Command, args []string) {
	used := []string{}
	if mif(cmd, "all") {
		used = append(used, "--all")
	}

	if mif(cmd, "up") {
		used = append(used, "--up")
	}

	if mif(cmd, "exited") {
		used = append(used, "--exited")
	}

	if len(args) > 0 {
		used = append(used, "INSTANCES")
	}

	if len(used) > 1 {
		logAndExit(log.WithField("flags", strings.Join(used, ",")), "Conflicting arguments provided")
	}
}

func getMultipleInstances(cmd *cobra.Command, args []string) []string {
	validateMultiInstanceFlags(cmd, args)

	if len(args) > 0 {
		// Instance names are case-insensitive on the command line but are always actually lower case
		lowerArgs := make([]string, len(args))
		for i, arg := range args {
			lowerArgs[i] = strings.ToLower(arg)
		}
		return lowerArgs
	}

	names := []string{}
	instances := app.GetInstances()
	for _, instance := range instances {
		if mif(cmd, "all") || (mif(cmd, "up") && strings.HasPrefix(instance.Status, "Up")) || (mif(cmd, "exited") && strings.HasPrefix(instance.Status, "Exited")) {
			names = append(names, instance.Name)
		}
	}

	return names
}

func addMultiInstanceFlags(cmd *cobra.Command, commandDesc string) {
	flags.AddFlag(cmd, "all", false, "Run "+commandDesc+" on all existing iscenv instances")
	flags.AddFlag(cmd, "up", false, "Run "+commandDesc+" on all running iscenv instances")
	flags.AddFlag(cmd, "exited", false, "Run "+commandDesc+" on all exited iscenv instances")
}

func mif(cmd *cobra.Command, name string) bool {
	return flags.GetBool(cmd, name)
}
