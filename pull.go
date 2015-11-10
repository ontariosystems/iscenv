/*
Copyright 2015 Ontario Systems

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

package main

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var pullCommand = &cobra.Command{
	Use:   "pull",
	Short: "Pull the latest ISC product versions",
	Long:  "Pull the latest versions of the ISC product images.",
}

func init() {
	pullCommand.Run = pull
}

func pull(_ *cobra.Command, _ []string) {
	imgopts := docker.PullImageOptions{Registry: REGISTRY, Repository: REPOSITORY, OutputStream: os.Stdout}
	authcfg := getAuthConfig()
	err := dockerClient.PullImage(imgopts, authcfg)
	if err != nil {
		fatalf("Could not pull latest ISC product version images, error: %s\n", err)
	}
}

func getAuthConfig() docker.AuthConfiguration {
	authcfg := docker.AuthConfiguration{}

	cfg, err := loadDefaultDockerConfig()
	if err != nil {
		fatalf("Error loading ~/.dockercfg, error: %s\n", err)
	}

	if cfg != nil {
		if auth, ok := cfg[REGISTRY]; ok {
			authcfg.Username, authcfg.Password, err = auth.credentials()
			if err != nil {
				fatalf("Could not parse credentials from ~/.dockercfg, error: %s\n", err)
			}
			authcfg.Email = auth.Email
		}
	}

	return authcfg
}
