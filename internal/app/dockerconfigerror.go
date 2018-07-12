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

// NewDockerConfigError creates and returns a new DockerConfigError
func NewDockerConfigError(path, registry string, err error) *DockerConfigError {
	return &DockerConfigError{
		ConfigPath: path,
		Registry:   registry,
		Err:        err,
	}
}

// DockerConfigError is an used for logging an error with a docker config
type DockerConfigError struct {
	ConfigPath string
	Registry   string
	Err        error
}

// Error returns the error string for the error associated with the DockerConfigError
func (dce *DockerConfigError) Error() string {
	return dce.Err.Error()
}
