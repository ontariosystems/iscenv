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
	"errors"
)

// Get the registry credentials from the default docker config
func GetRegistryCredentials(registry string) (username string, password string, email string, err error) {
	cfg, err := LoadDefaultDockerConfig()
	if err != nil {
		return "", "", "", NewDockerConfigError(cfg.Path, registry, err)
	}

	if cfg == nil {
		return "", "", "", NewDockerConfigError(cfg.Path, registry, errors.New("No docker configuration found"))
	}

	auth, ok := cfg.Entries[registry]
	if !ok {
		return "", "", "", NewDockerConfigError(cfg.Path, registry, errors.New("Registry not found in docker configuration"))
	}

	username, password, err = auth.Credentials()
	if err != nil {
		return "", "", "", NewDockerConfigError(cfg.Path, registry, errors.New("Could not retrieve credentials from docker configuration"))
	}

	return username, password, auth.Email, nil
}
