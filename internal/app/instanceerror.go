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

// NewInstanceError creates and returns a new InstanceError
func NewInstanceError(name, id string, err error) *InstanceError {
	return &InstanceError{
		InstanceName: name,
		InstanceID:   id,
		Err:          err,
	}
}

// InstanceError is an used for logging an error with an instance
type InstanceError struct {
	InstanceName string
	InstanceID   string
	Err          error
}

// Error returns the error string for the error associated with the InstanceError
func (ie *InstanceError) Error() string {
	return ie.Err.Error()
}
