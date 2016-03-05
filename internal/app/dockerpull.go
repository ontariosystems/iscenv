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
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

const (
	defaultRegistry = "https://index.docker.io/v1/"
)

// TODO: pass in the image
func DockerPull(image, tag string) error {

	authcfgs, err := docker.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		return err
	}

	registry := ""
	authcfg := docker.AuthConfiguration{}
	// If there's no . in the image that means it's the default registry (dots are only allowed in the host portion)
	if !strings.Contains(image, ".") {
		if cfg, ok := authcfgs.Configs[defaultRegistry]; ok {
			authcfg = cfg
		}
	} else {
		s := strings.Split(image, "/")
		registry = s[0]
		if cfg, ok := authcfgs.Configs[registry]; ok {
			authcfg = cfg
		} else if cfg, ok := authcfgs.Configs["https://"+registry]; ok {
			authcfg = cfg
		}
	}

	// TODO: when the image cna be passed in this will have to be parsed
	imgopts := docker.PullImageOptions{
		Registry:     registry,
		Repository:   image,
		Tag:          tag,
		OutputStream: ioutil.Discard, // TODO: Handle status updates... somehow
	}
	return DockerClient.PullImage(imgopts, authcfg)
}
