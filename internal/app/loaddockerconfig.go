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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/ontariosystems/iscenv/iscenv"
)

const NewDefaultDockerConfigName = ".docker/config.json"
const OldDefaultDockerConfigName = ".dockercfg"

// Will return nil, nil if the file simply doesn't exist
func LoadDefaultDockerConfig() (iscenv.DockerConfig, error) {
	cfgPath, err := FindDefaultDockerConfig()
	if err != nil {
		return iscenv.DockerConfig{}, err
	}

	return LoadDockerConfig(cfgPath)
}

func LoadDockerConfig(path string) (iscenv.DockerConfig, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := iscenv.DockerConfig{}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func FindDefaultDockerConfig() (string, error) {
	// TODO: Possibly use viper?
	env := os.Getenv("DOCKER_CONFIG")
	if env != "" {
		if FileExists(env) {
			return env, nil
		}
		return "", fmt.Errorf("DOCKER_CONFIG environment variable points to non-existent file, path: %s", env)
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	newCFGPath := path.Join(usr.HomeDir, NewDefaultDockerConfigName)
	if FileExists(newCFGPath) {
		return newCFGPath, nil
	}

	oldCFGPath := path.Join(usr.HomeDir, OldDefaultDockerConfigName)
	if FileExists(oldCFGPath) {
		return oldCFGPath, nil
	}

	return "", fmt.Errorf("Could not find Docker config at new or old default path, new: %s, old: %s", newCFGPath, oldCFGPath)
}
