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

package localversionsplugin

import (
	"os"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcuadros/go-version"
	"github.com/ontariosystems/iscenv/iscenv"
)

const (
	pluginKey = "local"
)

var dockerClient *docker.Client

type Plugin struct{}

func (plugin *Plugin) Main() {
	var err error
	if dockerClient, err = docker.NewClient(iscenv.DockerSocket); err != nil {
		log.WithError(err).Error("Failed to create docker client")
		os.Exit(1)
	}
	iscenv.ServeVersionsPlugin(plugin)
}

func (*Plugin) Key() string {
	return pluginKey
}

func (*Plugin) Versions(image string) (iscenv.ISCVersions, error) {
	localImages, err := dockerClient.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}

	vs := []string{}                        // for sorting
	vm := make(map[string]docker.APIImages) // for lookup
	for _, localImage := range localImages {
		for _, repoTag := range localImage.RepoTags {
			repo, tag := splitRepoTag(repoTag)
			if repo == image {
				if m, _ := regexp.MatchString("^\\d+(\\.\\d+)?", tag); m {
					vm[tag] = localImage
					vs = append(vs, tag)
				}
			}
		}
	}

	version.Sort(vs)

	versions := make(iscenv.ISCVersions, len(vs), len(vs))
	for i, v := range vs {
		ai := vm[v]
		versions[i] = &iscenv.ISCVersion{
			ID:      strings.TrimPrefix(ai.ID, "sha256:"),
			Version: v,
			Created: ai.Created,
			Source:  "local",
		}
	}

	return versions, nil
}

func splitRepoTag(repoTag string) (string, string) {
	s := strings.Split(repoTag, ":")
	return s[0], s[1]
}
