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
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/ontariosystems/isclib"
	log "github.com/Sirupsen/logrus"
)

type InstanceStateFn func(instance *isclib.Instance)

func NewISCInstanceManager(instanceName string, ccontrolPath string, csessionPath string, primaryCommand string, primaryCommandNamespace string) (*ISCInstanceManager, error) {

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
		Instance:                instance,
		PrimaryCommand:          primaryCommand,
		PrimaryCommandNamespace: primaryCommandNamespace,
	}

	return eim, nil
}

// Manages a instance within a container
type ISCInstanceManager struct {
	*isclib.Instance
	InstanceRunningHandler  InstanceStateFn
	PrimaryCommand          string
	PrimaryCommandNamespace string
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

	primchan := make(chan error)
	pclog := ilog.WithField("command", eim.PrimaryCommand).WithField("namespace", eim.PrimaryCommandNamespace)
	go eim.execPrimaryProcess(pclog, primchan)

	select {
	case sig := <-sigchan:
		log.WithField("signal", sig).Info("Received signal")
	case err := <-primchan:
		if err == nil {
			pclog.Info("Primary command complete")
		} else {
			pclog.WithError(err).Error("Error executing primary command")
		}
	}

	return eim.Stop()
}

func (eim *ISCInstanceManager) execPrimaryProcess(l log.FieldLogger, c chan<- error) {
	if eim.PrimaryCommand == "" {
		l.Debug("No primary command provided, skipping execute")
		return
	}

	l.Info("Starting primary command")

	cmd := eim.GetCSessionCommand(eim.PrimaryCommandNamespace, eim.PrimaryCommand)
	r, w := io.Pipe()
	cmd.Stdout = w

	var err error
	if err = cmd.Start(); err != nil {
		l.WithError(err).Error("Failed to start primary command csession")
		c <- err
		return
	}

	go func() {
		defer w.Close()
		err = cmd.Wait()
	}()

	RelogStream(l, true, r)

	c <- err
}
