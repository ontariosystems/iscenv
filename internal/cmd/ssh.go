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
	"strings"
	"unicode"

	"github.com/ontariosystems/iscenv/internal/iscenv"

	"github.com/spf13/cobra"
)

type sshExecFn func(string, []string) error

var sshFlags = struct {
	Command string
}{}

var sshCmd = &cobra.Command{
	Use:   "ssh INSTANCE",
	Short: "Connect to an instance",
	Long:  "This command is deprecated in favor of exec.  Connect to an instance with docker exec.  This command remains ssh for backwards compatibility but no longer actually uses ssh.",
	Run:   ssh,
}

func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringVarP(&sshFlags.Command, "command", "c", "", "Execute an SSH command directly")
}

func ssh(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		iscenv.Fatal("Must provide exactly 1 instance as the first argument")
	}

	var cmdArgs []string
	if sshFlags.Command != "" {
		cmdArgs = toArgs(sshFlags.Command)
	} else {
		cmdArgs = defaultExecCommand
	}

	err := iscenv.DockerExec(args[0], true, cmdArgs...)
	if err != nil {
		iscenv.Fatalf("Error running docker exec, error: %s", err)
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
