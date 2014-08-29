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

type tailExecFn func(string, []string) error

var tailExecFollow bool
var tailExecLines string
var tailExecFilename string

var tailCommand = &cobra.Command{
	Use:   "tail INSTANCE",
	Short: "tail a file within an instance",
	Long:  "Connect to an instance with SSH and tail the given file.",
}

var tailFilenames = map[string]string{
	"cconsole": "/ensemble/instances/docker/mgr/cconsole.log",
}

func init() {
	tailCommand.Run = tail
	tailCommand.Flags().BoolVarP(&tailExecFollow, "follow", "f", false, "Follow log output")
	tailCommand.Flags().StringVarP(&tailExecLines, "lines", "n", "K|all", "Output the last K lines; default is 10; or use -n +K to output lines starting with the Kth")
	tailCommand.Flags().StringVarP(&tailExecFilename, "file", "l", "cconsole", "Filename to tail. `cconsole` is a magic filename, and the default")
}

func tail(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		tailArgs := []string{"tail"}
		if tailExecFollow {
			tailArgs = append(tailArgs, "-f")
		}
		if tailExecLines != "" && tailExecLines != "K|all" {
			tailArgs = append(tailArgs, "-n", tailExecLines)
		}

		if magicFilename, ok := tailFilenames[tailExecFilename]; ok {
			tailExecFilename = magicFilename
		}
		if tailExecFilename == "" {
			tailArgs = append(tailArgs, tailFilenames["cconsole"])
		} else {
			tailArgs = append(tailArgs, tailExecFilename)
		}

		sshExec(args[0], nil, tailArgs...)
	} else {
		fatal("Must provide an instance")
	}
}
