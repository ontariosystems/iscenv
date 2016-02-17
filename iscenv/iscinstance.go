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
)

type ISCInstance struct {
	ID      string
	Name    string
	Version string
	Created int64
	Status  string
	Ports   ContainerPorts
}

func (i ISCInstance) PortOffset() (offset int64, err error) {
	var ss, web int64
	if ss, err = getOffset(i.Name, "SuperServer", i.Ports.SuperServer, PortExternalSS); err != nil {
		return -1, err
	}

	if web, err = getOffset(i.Name, "Web", i.Ports.Web, PortExternalWeb); err != nil {
		return -1, err
	}

	if web != ss {
		return -1, fmt.Errorf("Port offset for SuperServer and Web do not match, instance: %s, super server: %d, Web: %d", i.Name, ss, web)
	}

	return ss, nil
}

func getOffset(instanceName, portType string, port, basePort ContainerPort) (int64, error) {
	if port < basePort {
		return 0, fmt.Errorf("%s port is outside of range, instance: %s, port: %d, basePort: %d\n", portType, instanceName, port, basePort)
	}

	return int64(port - basePort), nil
}
