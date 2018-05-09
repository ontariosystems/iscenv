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
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
	"github.com/spf13/cobra"
)

var internalWaitForReadyCmd = &cobra.Command{
	Use:   "_waitforready",
	Short: "Waits for ISC product to be ready",
	Long:  "waits for an ISC product instance to be up and ready to use (this command is only available within containers)",
	Run:   internalWaitForReady,
}

func init() {
	if err := app.EnsureWithinContainer("_waitforready"); err != nil {
		return
	}

	rootCmd.AddCommand(internalWaitForReadyCmd)

	flags.AddConfigFlagP(internalWaitForReadyCmd, "timeout", "t", 600, "How long to wait, in seconds, for the instance to be ready before giving up")
}

func internalWaitForReady(cmd *cobra.Command, _ []string) {
	startTime := time.Now()

	var phase app.StartPhase
	for ; phase < app.StartPhaseInstanceRunning; phase = getStartPhase() {
		log.WithField("phase", phase).Debug("Not ready")
		time.Sleep(1 * time.Second)
		if int(time.Now().Sub(startTime).Seconds()) > flags.GetInt(cmd, "timeout") {
			logAndExit(log.StandardLogger(), "Instance failed to be ready in the allotted time")
		}
	}
	if phase != app.StartPhaseInstanceRunning {
		logAndExit(log.StandardLogger(), "Instance failed to be in ready state")
	}
}

func getStartPhase() app.StartPhase {
	resp, err := http.Get("http://localhost:59772")
	if err != nil {
		log.WithError(err).Warn("Failed to query health check service")
		return app.StartPhaseStartup
	}
	defer resp.Body.Close()

	status := app.NewStartStatus()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(status); err != nil {
		log.WithError(err).Warn("Failed to deserialize health check service response")
		return app.StartPhaseStartup
	}

	return status.Phase
}
