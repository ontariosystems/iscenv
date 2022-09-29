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

package app

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/v3/iscenv"
	log "github.com/sirupsen/logrus"
)

// DockerClient holds a reference to a client for interacting with docker
var DockerClient *docker.Client

func init() {
	var err error

	slog := log.WithField("dockerSocket", iscenv.DockerSocket)
	// Normally, I would not pre-emptively exit outside of the commands themselves but since this is an init there's not much choice
	if DockerClient, err = docker.NewClient(iscenv.DockerSocket); err != nil {
		ErrorLogger(slog, err).Panic("Failed to create docker client")
	}

	slog.Debug("Created docker client")
}
