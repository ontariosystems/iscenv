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
	"strconv"
)

// ContainerPorts is a listing of the ContainerPorts of the container
type ContainerPorts struct {
	SuperServer ContainerPort
	Web         ContainerPort
	HealthCheck ContainerPort
}

// ContainerPort is a int64 representation of a port in a container
type ContainerPort int64

// String returns a string representing the ContainerPort
func (p ContainerPort) String() string {
	return strconv.FormatInt(int64(p), 10)
}
