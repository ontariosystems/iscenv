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
	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
	"github.com/spf13/cobra"
)

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

	flags.AddFlagP(tailCmd, "follow", "f", false, "Follow log output")
	flags.AddFlagP(tailCmd, "lines", "n", "all", "Output all lines; or use -n K to output the last K lines; or use +K to output the Kth and following lines")
	flags.AddFlagP(tailCmd, "file", "l", "cconsole", "Filename to tail. `cconsole` is a magic filename, and the default")
}

func tail(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logAndExit(log.WithError(app.ErrSingleInstanceArg), "Invalid arguments")
	}

	instance, ilog := app.FindInstanceAndLogger(args[0])
	if instance == nil {
		logAndExit(ilog.WithError(app.ErrNoSuchInstance), "Invalid arguments")
	}

	if err := app.DockerExec(instance, false, buildTailArgs(cmd, args)...); err != nil {
		logAndExit(app.ErrorLogger(ilog.WithField("tailFile", flags.GetString(cmd, "file")), err), "Failed to tail file")
	}
}

func buildTailArgs(cmd *cobra.Command, args []string) []string {
	tailArgs := []string{"tail"}
	if flags.GetBool(cmd, "follow") {
		tailArgs = append(tailArgs, "-f")
	}

	lines := flags.GetString(cmd, "lines")
	if lines != "" {
		if lines == "all" {
			lines = "+0"
		}
		tailArgs = append(tailArgs, "-n", lines)
	}

	file := flags.GetString(cmd, "file")
	if actualFile, ok := tailFilenames[file]; ok {
		file = actualFile
	}

	if file == "" {
		tailArgs = append(tailArgs, tailFilenames["cconsole"])
	} else {
		tailArgs = append(tailArgs, file)
	}

	return tailArgs
}
