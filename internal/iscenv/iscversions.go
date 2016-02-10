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

package iscenv

import (
	"regexp"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	version "github.com/mcuadros/go-version"
)

type ISCVersions []*ISCVersion

func (evs ISCVersions) Latest() *ISCVersion {
	if len(evs) == 0 {
		Fatal("No ISC product versions exist run the iscenv pull command")
	}
	return evs[len(evs)-1]
}

func GetVersions() ISCVersions {
	images, err := DockerClient.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		Fatalf("Could not list images, error: %s\n", err)
	}

	vs := []string{}                        // for sorting
	vm := make(map[string]docker.APIImages) // for lookup
	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			repo, tag := splitRepoTag(repoTag)
			if repo == Repository {
				if m, _ := regexp.MatchString("^\\d+(\\.\\d+)?", tag); m {
					vm[tag] = image
					vs = append(vs, tag)
				}
			}
		}
	}

	version.Sort(vs)

	versions := make(ISCVersions, len(vs), len(vs))
	for i, v := range vs {
		ai := vm[v]
		versions[i] = &ISCVersion{ID: ai.ID, Version: v, Created: ai.Created}
	}

	return versions
}

func splitRepoTag(repoTag string) (string, string) {
	s := strings.Split(repoTag, ":")
	return s[0], s[1]
}
