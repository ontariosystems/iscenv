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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	log "github.com/sirupsen/logrus"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd"
	"github.com/ontariosystems/iscenv/iscenv"
)

func main() {
	exe, wrapped := iscenv.CalledAs()
	if wrapped {
		execInContainer(append([]string{exe}, os.Args[1:]...))
	} else {
		cmd.Execute()
	}
}

func execInContainer(args []string) {
	name := os.Getenv("ISCENV_INSTANCE")
	if name == "" {
		log.Panic("Must provide an instance in ISCENV_INSTANCE environment variable")
	}

	instance, _ := app.FindInstanceAndLogger(name)
	if name == "" {
		log.WithField("instance", name).Panic("Invalid instance")
	}

	// Since the whole point of faked executables is to trick wrappers, we need the
	// output to be untainted by additional logs.
	log.SetOutput(ioutil.Discard)

	var interactive bool
	switch strings.ToLower(os.Getenv("ISCENV_INTERACTIVE")) {
	case "true":
		interactive = true
	case "false":
		interactive = false
	default:
		interactive = terminal.IsTerminal(int(os.Stdin.Fd()))
	}

	if err := app.DockerExec(instance, interactive, args...); err != nil {
		if deerr, ok := err.(app.DockerExecError); ok {
			os.Exit(deerr.ExitCode)
		} else {
			fmt.Println(err)
			os.Exit(99)
		}
	}
}
