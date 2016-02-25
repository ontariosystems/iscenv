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
	"io/ioutil"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/iscenv"
)

// TODO: pass in the image
func DockerPull(version string) error {
	authcfg, err := GetAuthConfig(iscenv.Registry)
	if err != nil {
		return err
	}

	// TODO: when the image cna be passed in this will have to be parsed
	imgopts := docker.PullImageOptions{
		Registry:     iscenv.Registry,
		Repository:   iscenv.Repository,
		Tag:          version,
		OutputStream: ioutil.Discard, // TODO: Handle status updates... somehow
	}
	return DockerClient.PullImage(imgopts, authcfg)
}
