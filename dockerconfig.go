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
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

const defaultDockerConfigName = ".dockercfg"

type dockerConfigEntry struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

type dockerConfig map[string]dockerConfigEntry

func (dce dockerConfigEntry) credentials() (string, string, error) {
	creds, err := base64.StdEncoding.DecodeString(dce.Auth)
	if err != nil {
		return "", "", err
	}

	s := strings.Split(string(creds), ":")
	return s[0], s[1], nil
}

// Will return nil, nil if the file simply doesn't exist
func loadDefaultDockerConfig() (dockerConfig, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	cfgpath := path.Join(usr.HomeDir, defaultDockerConfigName)
	if _, err = os.Stat(cfgpath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return loadDockerConfig(cfgpath)
}

func loadDockerConfig(path string) (dockerConfig, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := dockerConfig{}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
