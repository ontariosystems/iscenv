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
	"strings"

	"github.com/ontariosystems/iscenv/v3/iscenv"
	log "github.com/sirupsen/logrus"
)

// FindInstance finds and returns the instance with the provided name
func FindInstance(instanceName string) *iscenv.ISCInstance {
	instanceName = strings.ToLower(instanceName)
	return GetInstances().Find(instanceName)
}

// FindInstanceAndLogger finds and returns the instance with the provided name and a logger for that instance
func FindInstanceAndLogger(instanceName string) (*iscenv.ISCInstance, *log.Entry) {
	var id string
	instance := FindInstance(instanceName)
	if instance != nil {
		id = instance.ID
	}

	return instance, InstanceLoggerArgs(instanceName, id)
}
