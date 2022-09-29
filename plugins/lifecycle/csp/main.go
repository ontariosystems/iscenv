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

package cspplugin

import (
	"encoding/json"
	"strconv"

	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/v3/internal/app"
	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
	log "github.com/sirupsen/logrus"
)

const (
	pluginKey       = "csp"
	defaultBasePort = 8443
)

var (
	plog = log.WithField("plugin", pluginKey)
)

// Plugin represents this plugin and serves as a place to attach functions to implement the Lifecycler interface
type Plugin struct{}

// Flags is used to hold flag values for this plugins flags
type Flags struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
	Port  int64  `json:"port"`
}

// Main serves as the main entry point for the plugin
func (plugin *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(plugin)
}

// Key returns the unique identifier for the plugin
func (*Plugin) Key() string {
	return pluginKey
}

// Flags returns an array of additional flags to add to the start command
func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	fb := iscenv.NewPluginFlagsBuilder()
	fb.AddFlag("image", true, "", "The docker image to use when creating the CSP container")
	fb.AddFlag("tag", true, "", "The tag of the docker image to use when creating the CSP container")
	fb.AddFlag("port", true, defaultBasePort, "The base port to use for CSP containers, the port offset will be added to this value to determine the actual port")
	return fb.Flags()
}

// Environment returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
func (*Plugin) Environment(version string, flags map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Copies returns an array of items to copy to the container in the format "src:dest"
func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Volumes returns an array of volumes to add where the string is a standard docker volume format "src:dest:flag"
func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Ports returns an array of additional ports to map in the format <optional hostIP>:hostPort:containerPort.  You may also prefix the host port with a + to indicate it should be shifted by the port offset
func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// AfterStart will run on the host after the container instance starts, receives the same flag values as start
func (plugin *Plugin) AfterStart(instance *iscenv.ISCInstance) error {
	flags, err := getFlagsForInstance(instance)
	if err != nil || flags.Image == "" {
		return err
	}

	l := plog.WithFields(log.Fields{"image": flags.Image, "tag": flags.Tag})
	l.Info("Downloading image for CSP container.  This may take a while.")
	if err := app.DockerPull(flags.Image, flags.Tag); err != nil {
		return err
	}

	po, err := instance.PortOffset()
	if err != nil {
		return err
	}

	// TODO: This is hacky.  The entire container handling portion of iscenv needs to be rewritten.  Right now it's more trouble than it's worth to restart this instance.  Too much of the port offset, etc. logic is rolled into simple container creation.  When it's rewritten remove this.
	if err := plugin.BeforeRemove(instance); err != nil {
		return err
	}

	name := getCSPContainerName(instance)
	l = l.WithField("name", name)
	l.Debug("Creating CSP container")
	id, err := app.DockerStart(app.DockerStartOptions{
		Name:                           name,
		FullName:                       name,
		Repository:                     flags.Image,
		Version:                        flags.Tag,
		PortOffset:                     po,
		PortOffsetSearch:               false,
		DisablePortOffsetConflictCheck: true,
		Entrypoint:                     nil,
		Command:                        nil,
		Environment:                    nil,
		Volumes:                        nil,
		Copies:                         nil,
		VolumesFrom:                    nil,
		ContainerLinks:                 []string{"iscenv-" + instance.Name + ":" + "iscenv"},
		Ports:                          []string{strconv.FormatInt(flags.Port+po, 10) + ":443"},
		Labels:                         nil,
		Recreate:                       false,
	})

	if err != nil {
		return err
	}

	l.WithField("id", id).Info("Created CSP container")
	return nil
}

// AfterStop will run on the host after the instance stops
func (*Plugin) AfterStop(instance *iscenv.ISCInstance) error {
	name := getCSPContainerName(instance)
	l := plog.WithField("name", name)
	if err := app.DockerClient.StopContainer(name, 30); err != nil {
		switch err := err.(type) {
		case *docker.NoSuchContainer:
			l.WithError(err).Debug("No CSP container to stop")
		default:
			return err
		}
	} else {
		l.Info("Stopped CSP container")
	}
	return nil
}

// BeforeRemove will run on the host before the instance is removed
func (*Plugin) BeforeRemove(instance *iscenv.ISCInstance) error {
	name := getCSPContainerName(instance)
	l := plog.WithField("name", name)
	if err := app.DockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: name, RemoveVolumes: true, Force: true}); err != nil {
		switch err := err.(type) {
		case *docker.NoSuchContainer:
			l.WithError(err).Debug("No CSP container to remove")
		default:
			return err
		}
	} else {
		l.Info("Removed CSP container")
	}
	return nil
}

// BeforeInstance will run within the container before the instance successfully starts
func (*Plugin) BeforeInstance(state *isclib.Instance) error {
	return nil
}

// WithInstance will run within the container after the instance starts
func (*Plugin) WithInstance(state *isclib.Instance) error {
	return nil
}

// AfterInstance will run within the container after the instance stops
func (*Plugin) AfterInstance(state *isclib.Instance) error {
	return nil
}

func getCSPContainerName(instance *iscenv.ISCInstance) string {
	return "csp-iscenv-" + instance.Name
}

func getFlagsForInstance(instance *iscenv.ISCInstance) (*Flags, error) {
	container, err := app.GetContainerForInstance(instance)
	if err != nil {
		return nil, err
	}

	flags := new(Flags)
	flagsStr, ok := container.Config.Labels["iscenv.lifecycler."+pluginKey+".flags"]
	if ok {
		if err := json.Unmarshal([]byte(flagsStr), flags); err != nil {
			return nil, err
		}

		if flags.Tag == "" {
			flags.Tag = "latest"
		}
	}

	return flags, nil
}
