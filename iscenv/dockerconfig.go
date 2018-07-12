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
	"encoding/base64"
	"strings"
)

// DockerConfig is the contents of a docker ocnfiguration file
type DockerConfig struct {
	Path    string
	Entries map[string]DockerConfigEntry
}

// DockerConfigEntry is a single entry from a docker configuration file
type DockerConfigEntry struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

// Credentials parses the credentials from this entry
func (dce DockerConfigEntry) Credentials() (user string, pass string, err error) {
	creds, err := base64.StdEncoding.DecodeString(dce.Auth)
	if err != nil {
		return "", "", err
	}

	s := strings.Split(string(creds), ":")
	return s[0], s[1], nil
}
