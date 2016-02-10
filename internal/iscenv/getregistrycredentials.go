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

// Get the registry credentials from the default docker config
func GetRegistryCredentials() (username string, password string, email string) {
	cfg, err := LoadDefaultDockerConfig()
	if err != nil {
		Fatalf("Error loading ~/.dockercfg, error: %s\n", err)
	}

	if cfg == nil {
		Fatalf("No ~/.dockercfg exists")
	}

	auth, ok := cfg[Registry]
	if !ok {
		Fatalf("No entry for %s in ~/.dockercfg", Registry)
	}

	username, password, err = auth.Credentials()
	if err != nil {
		Fatalf("Could not parse credentials from ~/.dockercfg, error: %s\n", err)
	}

	return username, password, auth.Email
}
