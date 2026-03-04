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

var (
	// ErrNoSuchInstance is an error used when the requested instance does not exist
	ErrNoSuchInstance = errors.New("no such instance")
	// ErrSingleInstanceArg is an error used when the wrong number of instances is provided
	ErrSingleInstanceArg = errors.New("must provide a single instance as the first argument")
	// ErrFailedToAddPluginFlags is an error used when adding flags for a plugin fails
	ErrFailedToAddPluginFlags = errors.New("failed to add flags from plugin")
	// ErrNotInContainer is an error used when the process is not executed in a container when it should be
	ErrNotInContainer = errors.New("not in Docker container")
	// ErrFailedEventPlugin is an error used when a plugin fails to execute properly
	ErrFailedEventPlugin = errors.New("failed to execute event plugin")
)
