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
	"github.com/spf13/cobra"
	"strings"
)

var stopTimeout uint

var stopCommand = &cobra.Command{
	Use:   "stop [OPTIONS] INSTANCE [INSTANCE...]",
	Short: "Stop an ISC product instance",
	Long:  "Stop a running ISC product instance, attempting a safe shutdown",
}

func init() {
	stopCommand.Run = stop
	stopCommand.Flags().UintVarP(&stopTimeout, "time", "t", 60, "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
}

func stop(_ *cobra.Command, args []string) {
	for _, arg := range args {
		instance := strings.ToLower(arg)
		current := getInstances()
		existing := current.find(instance)

		if existing != nil {
			err := dockerClient.StopContainer(existing.id, stopTimeout)
			if err != nil {
				Fatalf("Could not stop instance, name: %s, error: %s\n", existing.name, err)
			}

			fmt.Println(existing.id)
		} else {
			fmt.Printf("No such instance, name: %s\n", arg)
		}
	}
}
