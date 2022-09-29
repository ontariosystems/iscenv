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
	"time"

	"github.com/ontariosystems/iscenv/v3/iscenv"
)

// GetInstanceStartTime gets and returns the start time of the provided instance
func GetInstanceStartTime(instance *iscenv.ISCInstance) (time.Time, error) {
	var err error
	if container, err := GetContainerForInstance(instance); err == nil {
		return container.State.StartedAt, nil
	}

	return time.Time{}, err
}
