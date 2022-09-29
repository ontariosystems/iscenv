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
	"fmt"

	"github.com/ontariosystems/isclib/v2"
)

// The order here *must* match the actual order that things are started
const (
	StartPhaseStartup = iota
	StartPhaseInitPlugins
	StartPhaseInitManager
	StartPhaseEventBeforeInstance
	StartPhaseEventWithInstance
	StartPhaseInstanceRunning
	StartPhaseEventAfterInstance
	StartPhaseShutdown
)

// StartPhase represents a phase in the startup process
type StartPhase uint

// NewStartStatus creates and returns a new initialized StartStatus
func NewStartStatus() *StartStatus {
	return &StartStatus{
		Phase:         StartPhaseStartup,
		ActivePlugins: []string{},
		InstanceState: nil,
	}
}

// StartStatus represents information about the startup of an instance and the current phase it is in
type StartStatus struct {
	Phase           StartPhase       `json:"phase"`
	ActivePlugins   []string         `json:"activePlugins"`
	ExecutingPlugin string           `json:"executingPlugin"`
	InstanceState   *isclib.Instance `json:"instanceState"`
}

// Update will update the start phase of the provided instance to the provided start phase
func (ss *StartStatus) Update(phase StartPhase, state *isclib.Instance, executingPlugin string) {
	// Done this way rather than simply auto-advancing so the calling code is easier to read
	if ss.Phase != phase && ss.Phase+1 != phase {
		panic(fmt.Sprintf("Attempted to skip a phase or move backwards, current: %d, next: %d", ss.Phase, phase))
	}
	ss.Phase = phase
	if state != nil {
		ss.InstanceState = state
	}
	ss.ExecutingPlugin = executingPlugin
}
