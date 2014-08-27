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
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"strings"
)

var rmCommand = &cobra.Command{
	Use:   "rm INSTANCE [INSTANCE...]",
	Short: "Remove an ISC product instance",
	Long:  "Forcefully remove an ISC product instance regardless of its current state",
}

func init() {
	rmCommand.Run = rm
}

func rm(_ *cobra.Command, args []string) {
	for _, arg := range args {
		instance := strings.ToLower(arg)
		current := getInstances()
		existing := current.find(instance)

		if existing != nil {
			err := dockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: existing.id, RemoveVolumes: true, Force: true})
			if err != nil {
				Fatalf("Could not kill instance, name: %s, error: %s\n", existing.name, err)
			}
			fmt.Println(existing.id)
		} else {
			fmt.Printf("No such instance, name: %s\n", arg)
		}
	}
}
