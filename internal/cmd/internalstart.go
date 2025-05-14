/*
Copyright 2017 Ontario Systems

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
	"os"
	"path"

	"github.com/ontariosystems/iscenv/v3/internal/app"
	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	"github.com/ontariosystems/iscenv/v3/internal/plugins"
	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

	if err := addLifecyclerFlagsIfNeeded(internalStartCmd); err != nil {
		logAndExit(app.ErrorLogger(log.StandardLogger(), err), app.ErrFailedToAddPluginFlags.Error())
	}

	flags.AddConfigFlagP(internalStartCmd, "instance", "i", "", "The instance to manage")
	flags.AddConfigFlag(internalStartCmd, "ccontrol-path", "ccontrol", "The path to the ccontrol executable in the image")
	flags.AddConfigFlag(internalStartCmd, "csession-path", "csession", "The path to the csession executable in the image")
	addPrimaryCommandFlags(internalStartCmd)
}

func internalStart(cmd *cobra.Command, _ []string) {
	const zoneInfoPath string = "/usr/share/zoneinfo"
	const localTimePath string = "/etc/localtime"
	const timeZonePath string = "/etc/timezone"

	go startHealthCheck()

	if tz := os.Getenv("TZ"); tz != "" && os.Getuid() == 0 {
		log.WithField("time_zone", tz).Debug("Using provided time zone")
		if _, err := os.Stat(localTimePath); err == nil {
			if err := os.Remove(localTimePath); err != nil {
				logAndExit(app.ErrorLogger(log.StandardLogger().WithField("path", localTimePath), err), "Failed to remove previous time zone")
			}
			log.WithField("path", localTimePath).Debug("Removed previous time zone")
		}

		if err := os.Symlink(path.Join(zoneInfoPath, tz), localTimePath); err != nil {
			logAndExit(app.ErrorLogger(log.StandardLogger().WithField("path", localTimePath), err), "Failed to set time zone")
		}
		log.WithField("path", localTimePath).WithField("time_zone", tz).Debug("Set time zone")

		if err := os.WriteFile(timeZonePath, []byte(tz+"\n"), 0644); err != nil {
			logAndExit(app.ErrorLogger(log.StandardLogger().WithField("path", timeZonePath), err), "Failed to set time zone")
		}
		log.WithField("path", timeZonePath).WithField("time_zone", tz).Debug("Set time zone")
	}

	// We can't use the closing activator because we need the plugins to keep running the whole time that _start runs
	pluginsToActivate := getPluginsToActivate(rootCmd)
	startStatus.ActivePlugins = pluginsToActivate
	startStatus.Update(app.StartPhaseInitPlugins, nil, "")

	var lcs []*plugins.ActivatedLifecycler
	defer getActivatedLifecyclers(pluginsToActivate, getPluginArgs(), &lcs)(rootCtx)
	for _, lc := range lcs {
		startStatus.Update(app.StartPhaseInitPlugins, nil, lc.Id)
	}

	primaryCommand := flags.GetString(cmd, "primary-command")
	primaryCommandNS := flags.GetString(cmd, "primary-command-ns")
	if _, err := os.Stat(disablePrimaryCommandFile); err == nil {
		if err := os.Remove(disablePrimaryCommandFile); err != nil {
			logAndExit(app.ErrorLogger(log.StandardLogger(), err), "Failed to remove disable primary command file")
		}
		log.WithField("command", primaryCommand).WithField("namespace", primaryCommandNS).Warn("Disabling primary command; it will be re-enabled on next restart")
		primaryCommand = ""
	} else if !os.IsNotExist(err) {
		logAndExit(app.ErrorLogger(log.StandardLogger(), err), "Failed to check for existence of disable primary command file")
	}

	startStatus.Update(app.StartPhaseInitManager, nil, "")
	manager, err := app.NewISCInstanceManager(
		flags.GetString(cmd, "instance"),
		flags.GetString(cmd, "ccontrol-path"),
		flags.GetString(cmd, "csession-path"),
		primaryCommand,
		primaryCommandNS,
	)
	if err != nil {
		logAndExit(app.ErrorLogger(log.StandardLogger(), err), "Failed to create instance manager")
	}

	startStatus.Update(app.StartPhaseEventBeforeInstance, manager.Instance, "")
	for _, lc := range lcs {
		plog := lifecyclerLogger(lc, "BeforeInstance")
		plog.Info("Executing plugin")
		startStatus.Update(app.StartPhaseEventBeforeInstance, nil, lc.Id)
		if err := lc.Lifecycler.BeforeInstance(manager.Instance); err != nil {
			logAndExit(app.ErrorLogger(plog, err), app.ErrFailedEventPlugin.Error())
		}
	}

	manager.InstanceRunningHandler = func(*isclib.Instance) {
		startStatus.Update(app.StartPhaseEventWithInstance, manager.Instance, "")
		for _, lc := range lcs {
			plog := lifecyclerLogger(lc, "WithInstance")
			plog.Info("Executing plugin")
			startStatus.Update(app.StartPhaseEventWithInstance, nil, lc.Id)
			if err := lc.Lifecycler.WithInstance(manager.Instance); err != nil {
				logAndExit(app.ErrorLogger(plog, err), app.ErrFailedEventPlugin.Error())
			}
		}

		startStatus.Update(app.StartPhaseInstanceRunning, manager.Instance, "")
	}

	if err := manager.Manage(); err != nil {
		logAndExit(app.ErrorLogger(log.StandardLogger(), err), "Failed to manage instance")
	}

	startStatus.Update(app.StartPhaseEventAfterInstance, manager.Instance, "")
	for _, lc := range lcs {
		plog := lifecyclerLogger(lc, "AfterInstance")
		plog.Info("Executing plugin")
		startStatus.Update(app.StartPhaseEventAfterInstance, nil, lc.Id)
		if err := lc.Lifecycler.AfterInstance(manager.Instance); err != nil {
			logAndExit(app.ErrorLogger(plog, err), app.ErrFailedEventPlugin.Error())
		}
	}

	startStatus.Update(app.StartPhaseShutdown, nil, "")
	log.Info("Completed instance management")
}

func startHealthCheck() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(startStatus); err != nil {
			logAndExit(app.ErrorLogger(log.StandardLogger(), err), "Failed to encode JSON")
		}
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", iscenv.PortInternalHC), nil); err != nil && err != http.ErrServerClosed {
		logAndExit(app.ErrorLogger(log.StandardLogger(), err), "Health check listen failed")
	}
}

func lifecyclerLogger(lc *plugins.ActivatedLifecycler, method string) *log.Entry {
	return app.PluginLogger(lc.Id, method, lc.Path)
}
