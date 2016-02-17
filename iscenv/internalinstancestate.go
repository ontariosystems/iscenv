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
	"strconv"
	"strings"
)

type InternalInstanceState struct {
	Name            string
	Directory       string
	Version         string
	Status          InternalInstanceStatus
	Activity        string
	CPFFileName     string
	SuperServerPort int
	WebServerPort   int
	JDBCPort        int
	State           string
	// There appears to be an additional property after state but I don't know what it is!
}

func (iis *InternalInstanceState) Update(qlist string) (err error) {
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
