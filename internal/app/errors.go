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
	ErrNoSuchInstance         = errors.New("No such instance")
	ErrSingleInstanceArg      = errors.New("Must provide a single instance as the first argument")
	ErrFailedToAddPluginFlags = errors.New("Failed to add flags from plugin")
	ErrNotInContainer         = errors.New("Not in Docker container")
	ErrFailedEventPlugin      = errors.New("Failed to execute event plugin")
)
