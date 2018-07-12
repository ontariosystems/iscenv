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
	"fmt"
)

// ISCInstances is a slice of ISCInstance
type ISCInstances []*ISCInstance

// ByPortOffsets returns a map of ISCInstance indexed by their port offsets
func (is ISCInstances) ByPortOffsets() (map[int64]*ISCInstance, error) {
	offsets := make(map[int64]*ISCInstance)
	for _, i := range is {
		offset, err := i.PortOffset()
		if err != nil {
			return nil, err
		}
		offsets[offset] = i
	}

	return offsets, nil
}

// CalculatePortOffset returns the next available port offset
func (is ISCInstances) CalculatePortOffset(start int64) (int64, error) {
	offsets, err := is.ByPortOffsets()
	if err != nil {
		return -1, err
	}

	// The ports are spaced out by 1000 (56772, 57772) if there are more than 1k instances we'll collide.  1000 isc product instances on a single machine would be an exceptional number anyway.
	for i := start; i < 1000; i++ {
		if _, in := offsets[i]; !in {
			return i, nil
		}
	}

	return -1, fmt.Errorf("Could not determine next port offset")
}

// UsedPortOffset returns whether a port offset is currently in use
func (is ISCInstances) UsedPortOffset(offset int64) (bool, error) {
	offsets, err := is.ByPortOffsets()
	if err != nil {
		return false, err
	}

	_, used := offsets[offset]
	return used, nil
}

// Find returns the corresponding ISCInstance for the name string provided
func (is ISCInstances) Find(name string) *ISCInstance {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}

	return nil
}

// Exists returns whether an instance exists with the provided name
func (is ISCInstances) Exists(name string) bool {
	return is.Find(name) != nil
}
