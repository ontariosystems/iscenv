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
	"strings"

	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"

	"github.com/spf13/cobra"
)

type activateStarterFn func(id string, starter iscenv.Starter) error

var startFlags = struct {
	Plugins        string
	Remove         bool
	Version        string
	PortOffset     int64
	SearchForPort  bool
	Quiet          bool
	VolumesFrom    []string
	ContainerLinks []string
	CacheKeyURL    string
	PluginFlags    iscenv.PluginFlags
}{}

var startCmd = &cobra.Command{
	Use:   "start INSTANCE [INSTANCE...]",
	Short: "Start an ISC product container",
	Long:  "Create or start one or more ISC product containers with the supplied options",
	Run:   start,
}

func init() {
	// Since, we're adding flags and this has to happen in init, we're unfortunately going to have to load up and close the plugins here and in the start function, we could persist the manager globally but it's not as safe as a failure in init could concievably leave rpc servers running
	available := make([]string, 0)
	if err := activateStarters(nil, func(id string, starter iscenv.Starter) error {
		var err error
		available = append(available, id)
		startFlags.PluginFlags, err = starter.Flags()
		if err != nil {
			return fmt.Errorf("Could not retrieve plugin flags, plugin: %s, error: %s", id, err)
		}
		startFlags.PluginFlags.AddFlagsToFlagSet(id, startCmd.Flags())
		return nil
	}); err != nil {
		app.Fatalf("%s\n", err)
	}

	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVar(&startFlags.Plugins, "plugins", "", "An ordered comma-separated list of plugins you wish to activate. available plugins: "+strings.Join(available, ","))
	startCmd.Flags().BoolVarP(&startFlags.Remove, "rm", "", false, "Remove existing instance before starting.  By default, this switch will preserve port settings when recreating the instance.")
	startCmd.Flags().StringVarP(&startFlags.Version, "version", "v", "", "The version of ISC product to start.  By default this will find the latest version on your system.")
	startCmd.Flags().StringSliceVar(&startFlags.ContainerLinks, "link", nil, "Add link to another container.  They should be in the format 'iscenv-{iscenvname}', 'iscenv-{iscenvname}:{alias}' or '{containername}:{alias}'")
	startCmd.Flags().Int64VarP(&startFlags.PortOffset, "port-offset", "p", -1, "The port offset for this instance's ports.  -1 means autodetect.  Will increment by 1 if more than 1 instance is specified.")
	startCmd.Flags().BoolVarP(&startFlags.Quiet, "quiet", "q", false, "Only display numeric IDs")
	startCmd.Flags().StringSliceVar(&startFlags.VolumesFrom, "volumes-from", nil, "Mount volumes from the specified container(s)")
	startCmd.Flags().StringVarP(&startFlags.CacheKeyURL, "license-key-url", "k", "", "Download the cache.key file from the provided location rather than the default Statler URL")
	addMultiInstanceFlags(startCmd, "start")
}

func start(_ *cobra.Command, args []string) {
	pluginEnvironment, pluginVolumes, err := getPluginConfig(startFlags.PluginFlags, strings.Split(startFlags.Plugins, ","))
	if err != nil {
		app.Fatalf("Error loading environment and volumes from plugins, error: %s\n", err)
	}

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

	// TODO: Add the start command
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		id, err := app.DockerStart(app.DockerStartOptions{
			Name:             instance,
			Repository:       iscenv.Repository,
			Version:          startFlags.Version,
			PortOffset:       po,
			PortOffsetSearch: pos,
			Environment:      pluginEnvironment,
			Volumes:          pluginVolumes,
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

// if pluginsToActivate is nil (rather than an empty slice, it means all)
func activateStarters(pluginsToActivate []string, fn activateStarterFn) error {
	pm, err := app.NewPluginManager(iscenv.ApplicationName, iscenv.StarterKey, iscenv.StarterPlugin{})
	if err != nil {
		return err
	}
	defer pm.Close()

	if pluginsToActivate == nil {
		pluginsToActivate = pm.AvailablePlugins()
	}

	return pm.ActivatePlugins(pluginsToActivate, func(id string, raw interface{}) error {
		starter := raw.(iscenv.Starter)
		return fn(id, starter)
	})
}

func getPluginConfig(flags iscenv.PluginFlags, pluginsToActivate []string) (environment []string, volumes []string, err error) {
	environment = make([]string, 0)
	volumes = make([]string, 0)

	activateStarters(pluginsToActivate, func(id string, starter iscenv.Starter) error {
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
		return nil
	})

	return environment, volumes, nil
}
