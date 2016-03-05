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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
)

type starterInfo struct {
	ID      string
	Path    string
	Starter iscenv.Starter
}

var internalStartCmd = &cobra.Command{
	Use:   "_start",
	Short: "start/manage ISC product ",
	Long:  "manages an ISC product instance (this command is only available within containers)",
	Run:   internalStart,
}

var startStatus = app.NewStartStatus()

func init() {
	if err := app.EnsureWithinContainer("_start"); err != nil {
		return
	}

	rootCmd.AddCommand(internalStartCmd)

	if err := addStarterFlags(internalStartCmd); err != nil {
		app.ErrorLogger(nil, err).Fatal(app.ErrFailedToAddPluginFlags)
	}

	flags.AddConfigFlagP(internalStartCmd, "instance", "i", "", "The instance to manage")
	flags.AddConfigFlag(internalStartCmd, "ccontrol-path", "ccontrol", "The path to the ccontrol executable in the image")
	flags.AddConfigFlag(internalStartCmd, "csession-path", "csession", "The path to the csession executable in the image")
}

func internalStart(cmd *cobra.Command, _ []string) {

	go startHealthCheck()

	// We can't use the closing activator because we need the plugins to keep running the whole time that _start runs
	pluginsToActivate := strings.Split(flags.GetString(cmd, "plugins"), ",")
	startStatus.ActivePlugins = pluginsToActivate
	startStatus.Update(app.StartPhaseInitPlugins, nil, "")
	starters := make([]*starterInfo, len(pluginsToActivate))

	i := 0
	pm, err := activateStarters(
		pluginsToActivate,
		app.PluginArgs{
			LogLevel: flags.GetString(rootCmd, "log-level"),
			LogJSON:  flags.GetBool(rootCmd, "log-json"),
		},
		func(id, pluginPath string, starter iscenv.Starter) error {
			startStatus.Update(app.StartPhaseInitPlugins, nil, id)
			starters[i] = &starterInfo{
				ID:      id,
				Path:    pluginPath,
				Starter: starter,
			}
			i++
			return nil
		})

	if pm != nil {
		defer pm.Close()
	}

	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to activate plugins")
	}

	startStatus.Update(app.StartPhaseInitManager, nil, "")
	manager, err := app.NewInternalInstanceManager(
		flags.GetString(cmd, "instance"),
		flags.GetString(cmd, "ccontrol-path"),
		flags.GetString(cmd, "csession-path"),
	)
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to create instance manager")
	}

	startStatus.Update(app.StartPhaseEventBeforeInstance, manager.InternalInstance, "")
	for _, starter := range starters {
		plog := starterLogger(starter, "BeforeInstance")
		plog.Info("Executing plugin")
		startStatus.Update(app.StartPhaseEventBeforeInstance, nil, starter.ID)
		if err := starter.Starter.BeforeInstance(*manager.InternalInstance); err != nil {
			app.ErrorLogger(plog, err).Fatal(app.ErrFailedEventPlugin)
		}
	}

	manager.InstanceRunningHandler = func(*iscenv.InternalInstance) {
		startStatus.Update(app.StartPhaseEventWithInstance, manager.InternalInstance, "")
		for _, starter := range starters {
			plog := starterLogger(starter, "WithInstance")
			plog.Info("Executing plugin")
			startStatus.Update(app.StartPhaseEventWithInstance, nil, starter.ID)
			if err := starter.Starter.WithInstance(*manager.InternalInstance); err != nil {
				app.ErrorLogger(plog, err).Fatal(app.ErrFailedEventPlugin)
			}
		}

		startStatus.Update(app.StartPhaseInstanceRunning, manager.InternalInstance, "")
	}

	if err := manager.Manage(); err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to manage instance")
	}

	startStatus.Update(app.StartPhaseEventAfterInstance, manager.InternalInstance, "")
	for _, starter := range starters {
		plog := starterLogger(starter, "AfterInstance")
		plog.Info("Executing plugin")
		startStatus.Update(app.StartPhaseEventAfterInstance, nil, starter.ID)
		if err := starter.Starter.AfterInstance(*manager.InternalInstance); err != nil {
			app.ErrorLogger(plog, err).Fatal(app.ErrFailedEventPlugin)
		}
	}

	startStatus.Update(app.StartPhaseShutdown, nil, "")
}

func startHealthCheck() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(startStatus); err != nil {
			app.ErrorLogger(nil, err).Fatal("Failed to encode JSON")
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", iscenv.PortInternalHC), nil)
}

func starterLogger(si *starterInfo, method string) *log.Entry {
	return app.PluginLogger(si.ID, method, si.Path)
}
