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
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/iscenv/internal/cmd/flags"

	log "github.com/Sirupsen/logrus"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start INSTANCE [INSTANCE...]",
	Short: "Start an ISC product container",
	Long:  "Create or start one or more ISC product containers with the supplied options",
	Run:   start,
}

func init() {
	rootCmd.AddCommand(startCmd)
	if err := addStarterFlagsIfNotPluginCall(startCmd); err != nil {
		app.ErrorLogger(nil, err).Fatal(app.ErrFailedToAddPluginFlags)
	}

	addMultiInstanceFlags(startCmd, "start")
	flags.AddFlag(startCmd, "rm", false, "Remove existing instance before starting.  By default, this switch will preserve port settings when recreating the instance.")
	flags.AddFlagP(startCmd, "version", "v", "", "The version of ISC product to start.  By default this will find the latest version on your system.")
	flags.AddFlag(startCmd, "link", []string(nil), "Add link to another container.  They should be in the format 'iscenv-{iscenvname}', 'iscenv-{iscenvname}:{alias}' or '{containername}:{alias}'")
	flags.AddFlagP(startCmd, "port-offset", "p", int64(-1), "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
	flags.AddFlag(startCmd, "timeout", int64(300), "The number of seconds to wait on an instance to start before timing out.")
	flags.AddFlag(startCmd, "volumes-from", []string(nil), "Mount volumes from the specified container(s)")

	// Flags overriding the default settings *inside* of containers
	flags.AddConfigFlag(startCmd, "internal-instance", "docker", "The name of the actual ISC product instance within the container")
	flags.AddConfigFlag(startCmd, "superserver-port", int(iscenv.PortInternalSS), "The super server port inside the ISC product container")
	flags.AddConfigFlag(startCmd, "isc-http-port", int(iscenv.PortInternalWeb), "The ISC Web Server port inside the ISC product container")
	flags.AddConfigFlag(startCmd, "ccontrol-path", "ccontrol", "The path to the ccontrol executable within the container")
}

func start(cmd *cobra.Command, args []string) {
	log.Debug("Retrieving local versions")
	ensureImage()

	image := flags.GetString(rootCmd, "image")
	version := getVersion(image, flags.GetString(cmd, "version"))

	pluginsToActivate := strings.Split(flags.GetString(cmd, "plugins"), ",")
	environment, copies, volumes, ports, err := getPluginConfig(cmd, pluginsToActivate, version)
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to load container settings from plugin")
	}

	exe, err := osext.Executable()
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to determine iscenv executable path for bind mount")
	}

	// Add the iscenv executable itself as an item to copy into the container
	copies = append(copies, fmt.Sprintf("%s:%s", exe, iscenv.InternalISCEnvPath))

	// Add the standard ports

	ssPort := flags.GetInt(cmd, "superserver-port")
	httpPort := flags.GetInt(cmd, "isc-http-port")

	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalSS, ssPort))
	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalWeb, httpPort))
	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalHC, iscenv.PortInternalHC))

	// Add environment variables for the internal ports
	environment = append(environment, fmt.Sprintf("%s=%d", iscenv.EnvInternalSS, ssPort))
	environment = append(environment, fmt.Sprintf("%s=%d", iscenv.EnvInternalWeb, httpPort))
	environment = append(environment, fmt.Sprintf("%s=%d", iscenv.EnvInternalHC, iscenv.PortInternalHC))

	instances := getMultipleInstances(cmd, args)
	po := flags.GetInt64(cmd, "port-offset")
	pos := po < 0 || len(instances) > 1
	if po < 0 {
		po = 0
	}

	for _, instanceName := range instances {
		ilog := app.InstanceLoggerArgs(instanceName, "")
		ilog.Info("Starting instance")
		_, err := app.DockerStart(app.DockerStartOptions{
			Name:             instanceName,
			Repository:       image,
			Version:          version,
			PortOffset:       po,
			PortOffsetSearch: pos,
			Environment:      environment,
			Volumes:          volumes,
			Copies:           copies,
			Ports:            ports,
			Entrypoint:       []string{"/bin/iscenv", "_start"},
			Command: []string{
				// TODO: Plugin parameters and additional parameters passed from start itself (maybe)
				fmt.Sprintf("--instance=%s", flags.GetString(cmd, "internal-instance")),
				fmt.Sprintf("--ccontrol-path=%s", flags.GetString(cmd, "ccontrol-path")),
				fmt.Sprintf("--plugins=%s", flags.GetString(cmd, "plugins")),
				fmt.Sprintf("--log-level=%s", flags.GetString(rootCmd, "log-level")),
				"--log-json=true", // Always log JSON because we're going to proxy and parse it
			},
			VolumesFrom:    flags.GetStringSlice(cmd, "volumes-from"),
			ContainerLinks: flags.GetStringSlice(cmd, "link"),
			Recreate:       flags.GetBool(cmd, "rm"),
		})

		if err != nil {
			app.ErrorLogger(ilog, err).Fatal("Failed to create instance")
		}

		instance, ilog := app.FindInstanceAndLogger(instanceName)
		if instance == nil {
			ilog.Fatal("Failed to find newly created instance")
		}

		r, w := io.Pipe()
		go func() {
			if err := app.DockerLogs(instance, w); err != nil {
				app.ErrorLogger(ilog, err).Error("Error while retrieving container logs")
			}
		}()

		go func() {
			if err := app.WaitForInstance(instance, time.Duration(flags.GetInt64(cmd, "timeout"))*time.Second); err != nil {
				app.ErrorLogger(ilog, err).Fatal("Failed to start instance")
			}
			w.Close()
		}()

		decoder := json.NewDecoder(r)
		for {
			var rlog map[string]interface{}
			if err := decoder.Decode(&rlog); err == nil {
				relog(ilog, rlog)
			} else if err == io.EOF {
				break
			} else {
				app.ErrorLogger(ilog, err).Warn("Failed to parse remote json log")
			}
		}
		ilog.Info("Started instance")
	}
}

// getVersion will determine the appropriate version of the provided docker image to use and download it as needed.
// If the requestedVersion is not empty, it will use that version.
// If the requested version is empty, it will search for the latest local version for this supplied image.
// It returns the actual version to be used.
func getVersion(image, requestedVersion string) string {
	// Get the local versions (passing an empty plugins list means *only* local)
	versions, err := getVersions(image, []string{})
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to retrieve local versions")
	}

	// If no version was passed we will use the latest already downloaded image.  This gives some level of predictability to this feature.  You won't suddenly end up with a brand new field test version when you recreate an instance.
	version := requestedVersion
	if version == "" {
		if len(versions) == 0 {
			log.Fatal("No local versions from which to choose latest, must provide version with --version flag")
		}
		version = versions.Latest().Version
	}

	if !versions.Exists(version) {
		vlog := app.DockerRepoLogger(image).WithField("version", version)
		vlog.Info("Unable to find version locally, attempting to download.  This may take some time.")
		if err := app.DockerPull(image, version); err != nil {
			vlog.WithError(err).Fatal("Failed to pull image")
		}
	}

	return version
}

func getPluginConfig(cmd *cobra.Command, pluginsToActivate []string, version string) (environment, copies, volumes, ports []string, err error) {

	log.Debugf("Getting configuration from plugins: %v", len(pluginsToActivate))
	environment = make([]string, 0)
	copies = make([]string, 0)
	volumes = make([]string, 0)
	ports = make([]string, 0)

	if err := activateStartersAndClose(pluginsToActivate,
		app.PluginArgs{
			LogLevel: flags.GetString(rootCmd, "log-level"),
			LogJSON:  flags.GetBool(rootCmd, "log-json"),
		},
		func(id, pluginPath string, starter iscenv.Starter) error {
			flagValues := getPluginFlagValues(cmd, id)
			// Copy external plugin binaries to the /bin directory
			if pluginPath != "" {
				copies = append(copies, fmt.Sprintf("%s:%s/%s", pluginPath, iscenv.InternalISCEnvBinaryDir, filepath.Base(pluginPath)))
			}

			if env, err := starter.Environment(version, flagValues); err != nil {
				return app.NewPluginError(id, "Environment", pluginPath, err)
			} else if env != nil {
				environment = append(environment, env...)
			}

			if vols, err := starter.Volumes(version, flagValues); err != nil {
				return app.NewPluginError(id, "Volumes", pluginPath, err)
			} else if vols != nil {
				volumes = append(volumes, vols...)
			}

			if pts, err := starter.Ports(version, flagValues); err != nil {
				return app.NewPluginError(id, "Ports", pluginPath, err)
			} else if pts != nil {
				ports = append(ports, pts...)
			}
			return nil
		}); err != nil {
		return nil, nil, nil, nil, err
	}

	return environment, copies, volumes, ports, nil
}

func getPluginFlagValues(cmd *cobra.Command, plugin string) map[string]interface{} {
	flagValues := make(map[string]interface{})

	flog := log.WithField("plugin", plugin)
	prefix := cmd.Name() + "." + plugin + "-"
	for _, key := range flags.GetKeys() {
		if strings.HasPrefix(key, prefix) {
			value := flags.GetValueWithKey(key)
			flagValues[strings.TrimPrefix(key, prefix)] = value
			flog = flog.WithField(key, value)
		}
	}

	flog.Debug("Retrieved plugin flags")
	return flagValues
}

func relog(l *log.Entry, rlog map[string]interface{}) {
	for key, value := range rlog {
		if key != "level" && key != "msg" && key == "time" {
			l = l.WithField(key, value)
		}
	}

	level, ok := rlog["level"].(string)
	if !ok {
		l.WithField("level", rlog["level"]).Warn("Remote log level was not a string, using info")
		level = "info"
	}

	msg, ok := rlog["msg"].(string)
	if !ok {
		log.WithField("msg", msg).Error("Remote log message was not a string, skipping")
		return
	}

	switch level {
	case "debug":
		l.Debug(msg)
	case "info":
		l.Info(msg)
	case "warning":
		l.Warn(msg)
	case "error":
		l.Error(msg)
	default:
		l.WithField("origMsg", msg).Warn("Remote log with unknown level")
	}
}
