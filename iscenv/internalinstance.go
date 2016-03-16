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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	// This is the string which will be piped into a csession command to load the actual code to be executed into an in-memory buffer from a file.
	codeImportString = `try { ` +
		`set f="%s" ` +
		`open f:"R":1 ` +
		`if $test { use f zload  close f do MAIN halt } ` +
		`else { do $zutil(4, $job, 98) } } ` +
		`catch ex { ` +
		`do BACK^%%ETN ` +
		`use 0 ` +
		`write !,"Exception: ",ex.DisplayString(),!,` +
		`"  name: ",ex.Name,!,` +
		`"  code: ",ex.Code,! ` +
		`do $zutil(4, $job, 99) ` +
		`}`
)

type InternalInstance struct {
	// Required to be able to run the executor
	CSessionPath string `json:"-"`

	// These values come directly from ccontrol qlist
	Name            string                 `json:"name"`
	Directory       string                 `json:"directory"`
	Version         string                 `json:"version"`
	Status          InternalInstanceStatus `json:"status"`
	Activity        string                 `json:"activity"`
	CPFFileName     string                 `json:"cpfFileName"`
	SuperServerPort int                    `json:"superServerPort"`
	WebServerPort   int                    `json:"webServerPort"`
	JDBCPort        int                    `json:"jdbcPort"`
	State           string                 `json:"state"`
	// There appears to be an additional property after state but I don't know what it is!
}

func (iis *InternalInstance) Update(qlist string) (err error) {
	qs := strings.Split(qlist, "^")
	if len(qs) < 9 {
		return fmt.Errorf("Insufficient pieces in qlist, need at least 9, qlist: %s", qlist)
	}

	if iis.SuperServerPort, err = strconv.Atoi(qs[5]); err != nil {
		return err
	}

	if iis.WebServerPort, err = strconv.Atoi(qs[6]); err != nil {
		return err
	}

	if iis.JDBCPort, err = strconv.Atoi(qs[7]); err != nil {
		return err
	}

	iis.Name = qs[0]
	iis.Directory = qs[1]
	iis.Version = qs[2]
	iis.Status, iis.Activity = qlistStatus(qs[3])
	iis.CPFFileName = qs[4]
	iis.State = qs[8]

	return nil
}

func qlistStatus(statusAndTime string) (InternalInstanceStatus, string) {
	s := strings.SplitN(statusAndTime, ",", 2)
	var a string
	if len(s) > 1 {
		a = s[1]
	}
	return InternalInstanceStatus(strings.ToLower(s[0])), a
}

// This will execute the label MAIN from the properly formatted Cache INT code stored in the codeReader in namespace
func (ii *InternalInstance) Execute(namespace string, codeReader io.Reader) (output string, err error) {
	elog := log.WithField("namespace", namespace)

	codePath, err := ii.genExecutorTmpFile(codeReader)
	if err != nil {
		return "", err
	}

	defer os.Remove(codePath)

	cmd := exec.Command(ii.CSessionPath, ii.Name, "-U", namespace)

	in, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Start()
	importString := fmt.Sprintf(codeImportString, codePath)
	elog.WithField("importCode", importString).Debug("Attempting to load INT code into buffer")
	if _, err := in.Write([]byte(importString)); err != nil {
		return "", err
	}
	in.Close()

	elog.Debug("Waiting on csession to exit")
	err = cmd.Wait()
	return out.String(), err
}

func (ii *InternalInstance) genExecutorTmpFile(codeReader io.Reader) (string, error) {
	tmpFile, err := ioutil.TempFile("", "iscenv-ii-exec-")
	if err != nil {
		return "", err
	}

	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, codeReader); err != nil {
		return "", nil
	}

	return tmpFile.Name(), nil
}
