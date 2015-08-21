/*
Copyright 2015 Ontario Systems

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
	"strings"
)

var multiInstanceFlags multiInstanceOps

type multiInstanceOps struct {
	all    bool
	up     bool
	exited bool
}

func (this multiInstanceOps) validate(args []string) {
	used := []string{}
	if this.all {
		used = append(used, "--all")
	}

	if this.up {
		used = append(used, "--up")
	}

	if this.exited {
		used = append(used, "--exited")
	}

	if len(args) > 0 {
		used = append(used, "INSTANCES")
	}

	if len(used) > 1 {
		fatalf("Conflicting arguments provided: %s", strings.Join(used, ", "))
	}
}

func (this multiInstanceOps) getInstances(args []string) []string {
	this.validate(args)

	if len(args) > 0 {
		return args
	}

	names := []string{}
	instances := getInstances()
	for _, instance := range instances {
		if this.all || (this.up && strings.HasPrefix(instance.Status, "Up")) || (this.exited && strings.HasPrefix(instance.Status, "Exited")) {
			names = append(names, instance.Name)
		}
	}

	return names
}

func addMultiInstanceFlags(command *cobra.Command, commandDesc string) {
	command.Flags().BoolVarP(&multiInstanceFlags.all, "all", "", false, "Run "+commandDesc+" on all existing iscenv instances")
	command.Flags().BoolVarP(&multiInstanceFlags.up, "up", "", false, "Run "+commandDesc+" on all running iscenv instances")
	command.Flags().BoolVarP(&multiInstanceFlags.exited, "exited", "", false, "Run "+commandDesc+" on all exited iscenv instances")
}
