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

var tailFlags = struct {
	Follow   bool
	Lines    string
	Filename string
}{}

var tailCmd = &cobra.Command{
	Use:   "tail INSTANCE",
	Short: "tail a file within an instance; cconsole by default",
	Long:  "Connect to a container and tail the given file.",
	Run:   tail,
}

var tailFilenames = map[string]string{
	"cconsole": "/ensemble/instances/docker/mgr/cconsole.log",
}

func init() {
	rootCmd.AddCommand(tailCmd)

	tailCmd.Flags().BoolVarP(&tailFlags.Follow, "follow", "f", false, "Follow log output")
	tailCmd.Flags().StringVarP(&tailFlags.Lines, "lines", "n", "all", "Output all lines; or use -n K to output the last K lines; or use +K to output the Kth and following lines")
	tailCmd.Flags().StringVarP(&tailFlags.Filename, "file", "l", "cconsole", "Filename to tail. `cconsole` is a magic filename, and the default")
}

func tail(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		tailArgs := []string{"tail"}
		if tailFlags.Follow {
			tailArgs = append(tailArgs, "-f")
		}
		if tailFlags.Lines != "" {
			if tailFlags.Lines == "all" {
				tailFlags.Lines = "+0"
			}
			tailArgs = append(tailArgs, "-n", tailFlags.Lines)
		}

		if magicFilename, ok := tailFilenames[tailFlags.Filename]; ok {
			tailFlags.Filename = magicFilename
		}
		if tailFlags.Filename == "" {
			tailArgs = append(tailArgs, tailFilenames["cconsole"])
		} else {
			tailArgs = append(tailArgs, tailFlags.Filename)
		}

		iscenv.DockerExec(args[0], false, tailArgs...)
	} else {
		iscenv.Fatal("Must provide an instance")
	}
}
