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
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ontariosystems/iscenv/iscenv"
)

const (
	IIMDefaultCControlPath = "ccontrol"
)

type InstanceStateFn func(state *iscenv.InternalInstanceState)

func NewInternalInstanceManager(instanceName string, ccontrolPath string) (*InternalInstanceManager, error) {
	if ccontrolPath == "" {
		ccontrolPath = IIMDefaultCControlPath
	}

	iim := &InternalInstanceManager{
		InternalInstanceState: new(iscenv.InternalInstanceState),
		instanceName:          instanceName,
		ccontrolPath:          ccontrolPath,
	}

	if err := iim.Update(); err != nil {
		return nil, err
	}

	return iim, nil
}

// Manages a instance within a container
type InternalInstanceManager struct {
	instanceName string
	ccontrolPath string
	*iscenv.InternalInstanceState

	InstanceRunningHandler InstanceStateFn
}

func (iim *InternalInstanceManager) qlist() (string, error) {
	out, err := exec.Command(iim.ccontrolPath, "qlist", iim.instanceName).CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (iim *InternalInstanceManager) Manage() error {
	if err := iim.Start(); err != nil {
		return err
	}

	log.Printf("Started instance, name: %s, status: %s", iim.instanceName, iim.Status)
	if iim.InstanceRunningHandler != nil {
		iim.InstanceRunningHandler(iim.InternalInstanceState)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)

	// TODO: Deal with transient "run this during cache up" type commands, that should run and then cache should stop

	sig := <-sigchan
	log.Printf("Got signal: %s\n", sig)

	return iim.Stop()
}

// TODO: Think about a nozstu flag if there's a reason
func (iim *InternalInstanceManager) Start() error {
	if iim.Status.Down() {
		if output, err := exec.Command(iim.ccontrolPath, "start", iim.instanceName, "quietly").CombinedOutput(); err != nil {
			return fmt.Errorf("Error starting instance, error: %s, output: %s", err, output)
		}
	}

	if err := iim.Update(); err != nil {
		return fmt.Errorf("Error refreshing instance state during start, error: %s", err)
	}

	if !iim.Status.Running() {
		return fmt.Errorf("Failed to start instance, name: %s, status: %s", iim.instanceName, iim.Status)
	}

	return nil
}

func (iim *InternalInstanceManager) Stop() error {
	if iim.Status.Up() {
		args := []string{"stop", iim.instanceName}
		if iim.Status.RequiresBypass() {
			args = append(args, "bypass")
		}
		args = append(args, "quietly")
		if output, err := exec.Command(iim.ccontrolPath, args...).CombinedOutput(); err != nil {
			return fmt.Errorf("Error stopping instance, error: %s, output: %s", err, output)
		}
	}

	if err := iim.Update(); err != nil {
		return fmt.Errorf("Error refreshing instance state during stop, error: %s", err)
	}

	if !iim.Status.Down() {
		return fmt.Errorf("Failed to stop instance, name: %s, status: %s", iim.instanceName, iim.Status)
	}

	return nil
}

func (iim *InternalInstanceManager) Update() error {
	qlist, err := iim.qlist()
	if err != nil {
		return err
	}

	if err := iim.InternalInstanceState.Update(qlist); err != nil {
		return err
	}

	return nil
}
