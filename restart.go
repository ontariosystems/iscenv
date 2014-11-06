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

var restartCommand = &cobra.Command{
	Use:   "restart [OPTIONS] INSTANCE [INSTANCE...]",
	Short: "Restarts an ISC product instance",
	Long:  "Restart a running ISC product instance, attempting a safe shutdown",
}

func init() {
	restartCommand.Run = restart
	restartCommand.Flags().UintVarP(&stopTimeout, "time", "t", 60, "The amount of time to wait for the instance to stop cleanly before killing it.")
}

func restart(c *cobra.Command, args []string) {
	stop(c, args)
	start(c, args)
}
