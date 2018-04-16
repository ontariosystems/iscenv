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
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs INSTANCE",
	Short: "Fetch the logs of an instance",
	Long:  "Fetch the logs of an instance",
	Run:   displayLogs,
}

// TODO: Maybe add since flag
func init() {
	rootCmd.AddCommand(logsCmd)
	flags.AddFlagP(logsCmd, "follow", "f", false, "Follow log output")
	flags.AddFlag(logsCmd, "tail", "all", "Number of lines to show from the end of the logs")
}

func displayLogs(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logAndExit(log.WithError(app.ErrSingleInstanceArg), "Invalid arguments")
	}

	instance, ilog := app.FindInstanceAndLogger(args[0])
	if instance == nil {
		logAndExit(ilog.WithError(app.ErrNoSuchInstance), "Invalid arguments")
	}

	follow := flags.GetBool(cmd, "follow")
	since := int64(0)
	tail := flags.GetString(cmd, "tail")
	for {
		r, w := io.Pipe()

		go app.RelogStream(log.StandardLogger(), true, r)

		if err := app.DockerLogs(instance, since, tail, follow, w); err != nil {
			logAndExit(ilog.WithError(err), "Failed to retrieve docker logs")
		}

		since = time.Now().Unix()
		tail = "all"

		if !follow {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
}
