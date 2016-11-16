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
	"os"
	"strings"
	"unicode"

	log "github.com/Sirupsen/logrus"

	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh INSTANCE",
	Short: "Connect to an instance",
	Long:  "This command is deprecated in favor of exec.  Connect to an instance with docker exec.  This command remains ssh for backwards compatibility but no longer actually uses ssh.",
	Run:   ssh,
}

func init() {
	rootCmd.AddCommand(sshCmd)

	flags.AddFlagP(sshCmd, "command", "c", "", "Execute a command over SSH and exit")
}

func ssh(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal(app.ErrSingleInstanceArg)
	}

	instance, ilog := app.FindInstanceAndLogger(args[0])
	if instance == nil {
		ilog.Fatal(app.ErrNoSuchInstance)
	}

	var cmdArgs []string
	if command := flags.GetString(cmd, "command"); command != "" {
		cmdArgs = toArgs(command)
	} else {
		cmdArgs = defaultExecCommand
	}

	if err := app.DockerExec(instance, true, cmdArgs...); err != nil {
		if deerr, ok := err.(app.DockerExecError); ok {
			app.ErrorLogger(ilog, err).Error("Failed to run docker exec")
			os.Exit(deerr.ExitCode)
		} else {
			app.ErrorLogger(ilog, err).Fatal("Failed to run docker exec")
		}
	}
}

// If the arguments are too complicated, this will likely fall apart.  In that case, *DO NOT IMPROVE THIS* but point the user at exec
func toArgs(s string) []string {
	q := rune(0)
	f := func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	}

	return strings.FieldsFunc(s, f)
}
