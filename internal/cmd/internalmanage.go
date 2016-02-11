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

var internalManageFlags = struct {
	Instance     string
	CControlPath string
}{}

var internalManageCmd = &cobra.Command{
	Use:    "_manage",
	Short:  "internal: manage ISC product ",
	Long:   "DO NOT RUN THIS OUTSIDE OF AN INSTANCE CONTAINER. manages an ISC product instance",
	Hidden: true,
	Run:    internalManage,
}

func init() {
	rootCmd.AddCommand(internalManageCmd)

	internalManageCmd.Flags().StringVarP(&internalManageFlags.Instance, "instance", "i", "docker", "The instance to manage")
	internalManageCmd.Flags().StringVarP(&internalManageFlags.CControlPath, "ccontrolpath", "c", "ccontrol", "The path to the ccontrol executable in the image")
}

func internalManage(_ *cobra.Command, _ []string) {
	iscenv.EnsureWithinContainer("_manage")

	manager, err := iscenv.NewInternalInstanceManager(internalManageFlags.Instance, internalManageFlags.CControlPath)
	if err != nil {
		iscenv.Fatalf("Error creating instance manager, error: %s\n", err)
	}

	if err := manager.Manage(); err != nil {
		iscenv.Fatalf("Error managing instance, error: %s\n", err)
	}
}
