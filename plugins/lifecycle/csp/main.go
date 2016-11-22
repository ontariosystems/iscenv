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
	"os"
	"strconv"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/iscenv"
	"github.com/ontariosystems/iscenv/internal/app"
	"github.com/ontariosystems/isclib"
	log "github.com/Sirupsen/logrus"
)

const (
	pluginKey       = "csp"
	defaultBasePort = 8443
)

var (
	dockerClient *docker.Client
	plog         = log.WithField("plugin", pluginKey)
)

type Plugin struct{}
type Flags struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
	Port  int64  `json:"port"`
}

func (plugin *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(plugin)
}

func (*Plugin) Key() string {
	var err error
	if dockerClient, err = docker.NewClient(iscenv.DockerSocket); err != nil {
		plog.WithError(err).Error("Failed to create docker client")
		os.Exit(1)
	}
	return pluginKey
}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	fb := iscenv.NewPluginFlagsBuilder()
	fb.AddFlag("image", true, "", "The docker image to use when creating the CSP container")
	fb.AddFlag("tag", true, "", "The tag of the docker image to use when creating the CSP container")
	fb.AddFlag("port", true, defaultBasePort, "The base port to use for CSP containers, the port offset will be added to this value to determine the actual port")
	return fb.Flags()
}

func (*Plugin) Environment(version string, flags map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

func (p *Plugin) AfterStart(instance *iscenv.ISCInstance) error {
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

	// TODO: This is hacky.  The entire container handling portion of iscenv needs to be rewritten.  Right now it's more trouble that it's worth to restart this instance.  Too much of the port offset, etc. logic is rolled into simple container creation.  When it's rewritten remove this.
	if err := p.BeforeRemove(instance); err != nil {
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

func (*Plugin) BeforeInstance(state *isclib.Instance) error {
	return nil
}

func (*Plugin) WithInstance(state *isclib.Instance) error {
	return nil
}

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
