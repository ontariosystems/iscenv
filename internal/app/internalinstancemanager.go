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
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/iscenv/iscenv"
)

const (
	IIMDefaultCControlPath = "ccontrol"
	IIMDefaultCSessionPath = "csession"
)

type InstanceStateFn func(state *iscenv.InternalInstance)

func NewInternalInstanceManager(instanceName string, ccontrolPath string, csessionPath string) (*InternalInstanceManager, error) {
	if ccontrolPath == "" {
		ccontrolPath = IIMDefaultCControlPath
	}

	iim := &InternalInstanceManager{
		InternalInstance: &iscenv.InternalInstance{CSessionPath: csessionPath},
		instanceName:     instanceName,
		ccontrolPath:     ccontrolPath,
		csessionPath:     csessionPath,
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
	csessionPath string
	*iscenv.InternalInstance

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
	ilog := log.WithField("name", iim.instanceName)
	ilog.Debug("Starting instance")
	if err := iim.Start(); err != nil {
		return err
	}

	if iim.InstanceRunningHandler != nil {
		ilog.Debug("Executing instance running handler")
		iim.InstanceRunningHandler(iim.InternalInstance)
	}

	ilog.WithField("status", iim.Status).Info("Started instance")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)

	// TODO: Add a stop immediately flag that allows you to just run the instance running handler and then exit

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
	ilog := log.WithField("name", iim.instanceName)
	ilog.Debug("Shutting down instance")
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

	if err := iim.InternalInstance.Update(qlist); err != nil {
		return err
	}

	return nil
}
