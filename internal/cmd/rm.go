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
	"fmt"
	"strings"

	"github.com/ontariosystems/iscenv/internal/app"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm INSTANCE [INSTANCE...]",
	Short: "Remove an ISC product instance",
	Long:  "Forcefully remove an ISC product instance regardless of its current state",
	Run:   rm,
}

func init() {
	rootCmd.AddCommand(rmCmd)

	addMultiInstanceFlags(rmCmd, "rm")
}

func rm(_ *cobra.Command, args []string) {
	instances := multiInstanceFlags.getInstances(args)
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		current := app.GetInstances()
		existing := current.Find(instance)

		if existing != nil {
			err := app.DockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: existing.ID, RemoveVolumes: true, Force: true})
			if err != nil {
				app.Fatalf("Could not kill instance, name: %s, error: %s\n", existing.Name, err)
			}
			fmt.Println(existing.ID)
		} else {
			fmt.Printf("No such instance, name: %s\n", instanceName)
		}
	}
}
