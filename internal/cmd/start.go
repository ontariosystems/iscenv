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
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"

	log "github.com/Sirupsen/logrus"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var startFlags = struct {
	Remove         bool
	Version        string
	PortOffset     int64
	VolumesFrom    []string
	ContainerLinks []string
	StartTimeout   int64
	Plugins        string
	PluginFlags    map[string]*iscenv.PluginFlags
}{}

var startCmd = &cobra.Command{
	Use:   "start INSTANCE [INSTANCE...]",
	Short: "Start an ISC product container",
	Long:  "Create or start one or more ISC product containers with the supplied options",
	Run:   start,
}

func init() {
	// Since, we're adding flags and this has to happen in init, we're unfortunately going to have to load up and close the plugins here and in the start function, we could persist the manager globally but it's not as safe as a failure in init could concievably leave rpc servers running
	startFlags.PluginFlags = make(map[string]*iscenv.PluginFlags)
	if err := addStarterFlags(startCmd, &startFlags.Plugins, startFlags.PluginFlags); err != nil {
		app.ErrorLogger(nil, err).Fatal(app.ErrFailedToAddPluginFlags)
	}

	addMultiInstanceFlags(startCmd, "start")
	startCmd.Flags().BoolVarP(&startFlags.Remove, "rm", "", false, "Remove existing instance before starting.  By default, this switch will preserve port settings when recreating the instance.")
	startCmd.Flags().StringVarP(&startFlags.Version, "version", "v", "", "The version of ISC product to start.  By default this will find the latest version on your system.")
	startCmd.Flags().StringSliceVar(&startFlags.ContainerLinks, "link", nil, "Add link to another container.  They should be in the format 'iscenv-{iscenvname}', 'iscenv-{iscenvname}:{alias}' or '{containername}:{alias}'")
	startCmd.Flags().Int64VarP(&startFlags.PortOffset, "port-offset", "p", -1, "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
	startCmd.Flags().Int64Var(&startFlags.StartTimeout, "start-timeout", 300, "The number of seconds to wait on an instance to start before timing out.")
	startCmd.Flags().StringSliceVar(&startFlags.VolumesFrom, "volumes-from", nil, "Mount volumes from the specified container(s)")

	rootCmd.AddCommand(startCmd)
}

func start(_ *cobra.Command, args []string) {
	environment, volumes, ports, err := getPluginConfig(startFlags.Version, startFlags.PluginFlags, strings.Split(startFlags.Plugins, ","))
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to load container settings from plugin")
	}

	exe, err := osext.Executable()
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to determine iscenv executable path for bind mount")
	}

	// Add the iscenv executable itself as a volume
	volumes = append(volumes, fmt.Sprintf("%s:%s:ro", exe, iscenv.InternalISCEnvPath))

	// Add the standard ports
	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalSS, iscenv.PortInternalSS))
	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalWeb, iscenv.PortInternalWeb))
	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalHC, iscenv.PortInternalHC))

	rlog := app.DockerRepoLogger(iscenv.Repository)
	log.Debug("Retrieving local versions")
	// Get the local versions (passing an empty plugins list means *only* local)
	versions, err := getVersions(iscenv.Repository, []string{})
	if err != nil {
		app.ErrorLogger(rlog, err).Fatal("Failed to retrieve versions")
	}

	// If no version was passed we will use the latest already downloaded image.  This gives some level of predictability to this feature.  You won't suddenly end up with a brand new field test version when you recreate an instance.
	if startFlags.Version == "" {
		if len(versions) == 0 {
			rlog.Fatal("No local versions, must provide version with --version flag")
		}
		startFlags.Version = versions.Latest().Version
	} else {
		if !versions.Exists(startFlags.Version) {
			vlog := rlog.WithField("version", startFlags.Version)
			vlog.Info("Unable to find version locally, attempting to download.  This may take some time.")
			if err := app.DockerPull(startFlags.Version); err != nil {
				vlog.WithError(err).Fatal("Failed to pull image")
			}
		}
	}

	instances := multiInstanceFlags.getInstances(args)

	po := startFlags.PortOffset
	pos := po < 0 || len(instances) > 1
	if po < 0 {
		po = 0
	}

	for _, instanceName := range instances {
		ilog := app.InstanceLoggerArgs(instanceName, "")
		ilog.Info("Starting instance")
		_, err := app.DockerStart(app.DockerStartOptions{
			Name:             instanceName,
			Repository:       iscenv.Repository,
			Version:          startFlags.Version,
			PortOffset:       po,
			PortOffsetSearch: pos,
			Environment:      environment,
			Volumes:          volumes,
			Ports:            ports,
			Entrypoint:       []string{"/bin/iscenv", "_start"},
			Command: []string{
				// TODO: Plugin parameters and additional parameters passed from start itself (maybe)
				fmt.Sprintf("--plugins=%s", startFlags.Plugins),
				fmt.Sprintf("--log-level=%s", globalFlags.LogLevel),
				fmt.Sprintf("--log-json=%t", globalFlags.LogJSON),
			},
			VolumesFrom:    startFlags.VolumesFrom,
			ContainerLinks: startFlags.ContainerLinks,
			Recreate:       startFlags.Remove,
		})

		if err != nil {
			app.ErrorLogger(ilog, err).Fatal("Failed to create instance")
		}

		instance, ilog := app.FindInstanceAndLogger(instanceName)
		if instance == nil {
			ilog.Fatal("Failed to find newly created instance")
		}

		if err := app.WaitForInstance(instance, time.Duration(startFlags.StartTimeout)*time.Second); err != nil {
			app.ErrorLogger(ilog, err).Fatal("Failed to start instance")
		}

		ilog.Info("Started instance")
	}
}

func getPluginConfig(version string, pluginFlags map[string]*iscenv.PluginFlags, pluginsToActivate []string) (environment []string, volumes []string, ports []string, err error) {

	environment = make([]string, 0)
	volumes = make([]string, 0)
	ports = make([]string, 0)

	if err := activateStartersAndClose(pluginsToActivate,
		app.PluginArgs{
			LogLevel: globalFlags.LogLevel,
			LogJSON:  globalFlags.LogJSON,
		},
		func(id, pluginPath string, starter iscenv.Starter) error {
			// Mount the plugin itself into the /bin directory
			volumes = append(volumes, fmt.Sprintf("%s:%s/%s:ro", pluginPath, iscenv.InternalISCEnvBinaryDir, filepath.Base(pluginPath)))
			if env, err := starter.Environment(version, *pluginFlags[id]); err != nil {
				return app.NewPluginError(id, "Environment", pluginPath, err)
			} else if env != nil {
				environment = append(environment, env...)
			}

			if vols, err := starter.Volumes(version, *pluginFlags[id]); err != nil {
				return app.NewPluginError(id, "Volumes", pluginPath, err)
			} else if vols != nil {
				volumes = append(volumes, vols...)
			}

			if pts, err := starter.Ports(version, *pluginFlags[id]); err != nil {
				return app.NewPluginError(id, "Ports", pluginPath, err)
			} else if pts != nil {
				ports = append(ports, pts...)
			}
			return nil
		}); err != nil {
		return nil, nil, nil, err
	}

	return environment, volumes, ports, nil
}
