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
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var startFlags = struct {
	Remove         bool
	Version        string
	PortOffset     int64
	SearchForPort  bool
	Quiet          bool
	VolumesFrom    []string
	ContainerLinks []string
	CacheKeyURL    string
	Plugins        string
	PluginFlags    iscenv.PluginFlags
}{}

var startCmd = &cobra.Command{
	Use:   "start INSTANCE [INSTANCE...]",
	Short: "Start an ISC product container",
	Long:  "Create or start one or more ISC product containers with the supplied options",
	Run:   start,
}

func init() {
	log.SetOutput(ioutil.Discard) // This is to silence the logging from go-plugin

	// Since, we're adding flags and this has to happen in init, we're unfortunately going to have to load up and close the plugins here and in the start function, we could persist the manager globally but it's not as safe as a failure in init could concievably leave rpc servers running
	if err := addStarterFlags(startCmd, &startFlags.Plugins, &startFlags.PluginFlags); err != nil {
		app.Fatalf("%s\n", err)
	}

	addMultiInstanceFlags(startCmd, "start")
	startCmd.Flags().BoolVarP(&startFlags.Remove, "rm", "", false, "Remove existing instance before starting.  By default, this switch will preserve port settings when recreating the instance.")
	startCmd.Flags().StringVarP(&startFlags.Version, "version", "v", "", "The version of ISC product to start.  By default this will find the latest version on your system.")
	startCmd.Flags().StringSliceVar(&startFlags.ContainerLinks, "link", nil, "Add link to another container.  They should be in the format 'iscenv-{iscenvname}', 'iscenv-{iscenvname}:{alias}' or '{containername}:{alias}'")
	startCmd.Flags().Int64VarP(&startFlags.PortOffset, "port-offset", "p", -1, "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
	startCmd.Flags().BoolVarP(&startFlags.Quiet, "quiet", "q", false, "Only display numeric IDs")
	startCmd.Flags().StringSliceVar(&startFlags.VolumesFrom, "volumes-from", nil, "Mount volumes from the specified container(s)")
	startCmd.Flags().StringVarP(&startFlags.CacheKeyURL, "license-key-url", "k", "", "Download the cache.key file from the provided location rather than the default Statler URL")

	rootCmd.AddCommand(startCmd)
}

func start(_ *cobra.Command, args []string) {
	environment, volumes, ports, err := getPluginConfig(startFlags.PluginFlags, strings.Split(startFlags.Plugins, ","))
	if err != nil {
		app.Fatalf("Error loading environment and volumes from plugins, error: %s\n", err)
	}

	exe, err := osext.Executable()
	if err != nil {
		app.Fatalf("Could not determine iscenv executable path for bind mount")
	}

	// Add the iscenv executable itself as a volume
	volumes = append(volumes, fmt.Sprintf("%s:%s:ro", exe, iscenv.InternalISCEnvPath))

	// Add the standard ports
	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalSS, iscenv.PortInternalSS))
	ports = append(ports, fmt.Sprintf("+%d:%d", iscenv.PortExternalWeb, iscenv.PortInternalWeb))

	// TODO: latest should only get the latest local version when the version plugins exist
	if startFlags.Version == "" {
		startFlags.Version = app.GetVersions().Latest().Version
	}

	instances := multiInstanceFlags.getInstances(args)

	po := startFlags.PortOffset
	pos := po < 0 || len(instances) > 1
	if po < 0 {
		po = 0
	}

	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		id, err := app.DockerStart(app.DockerStartOptions{
			Name:             instance,
			Repository:       iscenv.Repository,
			Version:          startFlags.Version,
			PortOffset:       po,
			PortOffsetSearch: pos,
			Environment:      environment,
			Volumes:          volumes,
			Ports:            ports,
			Entrypoint:       []string{"/bin/iscenv", "_start"},
			Command:          []string{fmt.Sprintf("--plugins=%s", startFlags.Plugins)}, // TODO: Plugin parameters and parameters passed from start itself
			VolumesFrom:      startFlags.VolumesFrom,
			ContainerLinks:   startFlags.ContainerLinks,
			Recreate:         startFlags.Remove,
		})
		if err != nil {
			app.Fatalf("Could not create instance, name: %s, error: %s", instance, err)
		}
		fmt.Println(id)
	}
}

func getPluginConfig(flags iscenv.PluginFlags, pluginsToActivate []string) (environment []string, volumes []string, ports []string, err error) {

	environment = make([]string, 0)
	volumes = make([]string, 0)
	ports = make([]string, 0)

	if err := activateStartersAndClose(pluginsToActivate, func(id, pluginPath string, starter iscenv.Starter) error {
		// Mount the plugin itself into the /bin directory
		volumes = append(volumes, fmt.Sprintf("%s:%s/%s:ro", pluginPath, iscenv.InternalISCEnvBinaryDir, filepath.Base(pluginPath)))
		if env, err := starter.Environment(startFlags.PluginFlags); err != nil {
			return err
		} else if env != nil {
			environment = append(environment, env...)
		}

		if vols, err := starter.Volumes(startFlags.PluginFlags); err != nil {
			return err
		} else if vols != nil {
			volumes = append(volumes, vols...)
		}

		if pts, err := starter.Ports(startFlags.PluginFlags); err != nil {
			return err
		} else if pts != nil {
			ports = append(ports, pts...)
		}
		return nil
	}); err != nil {
		return nil, nil, nil, err
	}

	return environment, volumes, ports, nil
}
