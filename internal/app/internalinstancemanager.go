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
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/isclib"
)

type InstanceStateFn func(instance *isclib.Instance)

func NewISCInstanceManager(instanceName string, ccontrolPath string, csessionPath string) (*ISCInstanceManager, error) {

	if ccontrolPath != "" {
		isclib.SetCControlPath(ccontrolPath)
	}

	if csessionPath != "" {
		isclib.SetCSessionPath(csessionPath)
	}

	instance, err := isclib.LoadInstance(instanceName)
	if err != nil {
		return nil, err
	}

	eim := &ISCInstanceManager{
		Instance: instance,
	}

	return eim, nil
}

// Manages a instance within a container
type ISCInstanceManager struct {
	*isclib.Instance
	InstanceRunningHandler InstanceStateFn
}

func (eim *ISCInstanceManager) Manage() error {
	ilog := log.WithField("name", eim.Instance.Name)
	ilog.Debug("Starting instance")
	if err := eim.Instance.Start(); err != nil {
		return err
	}

	if eim.InstanceRunningHandler != nil {
		ilog.Debug("Executing instance running handler")
		eim.InstanceRunningHandler(eim.Instance)
	}

	ilog.WithField("status", eim.Status).Info("Started instance")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)

	// TODO: Add a stop immediately flag that allows you to just run the instance running handler and then exit

	sig := <-sigchan
	log.Printf("Got signal: %s\n", sig)

	return eim.Stop()
}
