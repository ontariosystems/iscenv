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

type ISCInstances []*ISCInstance

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

func (is ISCInstances) CalculatePortOffset(start int64) (int64, error) {
	offsets, err := is.ByPortOffsets()
	if err != nil {
		return -1, err
	}

	for i := start; i < 65535; i++ {
		if _, in := offsets[i]; !in {
			return i, nil
		}
	}

	return -1, fmt.Errorf("Could not determine next port offset")
}

func (is ISCInstances) UsedPortOffset(offset int64) (bool, error) {
	offsets, err := is.ByPortOffsets()
	if err != nil {
		return false, err
	}

	_, used := offsets[offset]
	return used, nil
}

func (is ISCInstances) Find(name string) *ISCInstance {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}

	return nil
}

func (is ISCInstances) Exists(name string) bool {
	return is.Find(name) != nil
}
