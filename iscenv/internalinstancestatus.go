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

type InternalInstanceStatus string

const (
	InstanceStateUnknown           = ""
	InstanceStateRunning           = "running"
	InstanceStateInhibited         = "sign-on inhibited"
	InstanceStatePrimaryTransition = "sign-on inhibited:primary transition"
	InstanceStateDown              = "down"
	InstanceStateMissingIDS        = "running on node ? (cache.ids missing)"
)

// Returns true when the status is known and can be handled by this code
func (iis InternalInstanceStatus) Handled() bool {
	switch iis {
	default:
		return false
	case
		InstanceStateRunning,
		InstanceStateInhibited,
		InstanceStatePrimaryTransition,
		InstanceStateDown,
		InstanceStateMissingIDS:
		return true
	}
}

// Returns true if Cache is running and could be used normally
func (iis InternalInstanceStatus) Running() bool {
	switch iis {
	default:
		return false
	case
		InstanceStateRunning,
		InstanceStateMissingIDS:
		return true
	}
}

// Returns true if Cache is in a state where it is up but not necessarily cleanly
func (iis InternalInstanceStatus) Up() bool {
	switch iis {
	default:
		return false
	case
		InstanceStateRunning,
		InstanceStateInhibited,
		InstanceStateMissingIDS:
		return true
	}
}

// Returns true when the instance is down
func (iis InternalInstanceStatus) Down() bool {
	switch iis {
	default:
		return false
	case
		InstanceStateDown:
		return true
	}
}

// Returns true when a bypass is required to stop the instance
func (iis InternalInstanceStatus) RequiresBypass() bool {
	switch iis {
	default:
		return false
	case
		InstanceStateInhibited,
		InstanceStatePrimaryTransition:
		return true
	}
}
