/*
Copyright 2017 Ontario Systems

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
	"strings"

	"github.com/ontariosystems/iscenv/v3/iscenv"
	log "github.com/sirupsen/logrus"
)

// EnsureWithinContainer returns an error if it is executed from outside a container
func EnsureWithinContainer(commandName string) error {
	// if iscenv started the container, this will be set
	if _, ok := os.LookupEnv(iscenv.EnvInternalContainer); ok {
		return nil
	}

	proc1CGroupContents, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		e := fmt.Errorf("Failed to determine environment")
		log.WithField("path", "/proc/1/cgroup").WithError(err).Error(e)
		return e
	}

	// if we have some control groups owned by docker, then we are within a container
	contents := string(proc1CGroupContents)
	cgroup_substrs := []string{
		":/docker/",
		":/kubepods/",
		":/kubepods.slice/",
		":/system.slice/docker-",
		":/system.slice/system.slice:docker:",
	}
	for _, substr := range cgroup_substrs {
		if strings.Contains(contents, substr) {
			return nil
		}
	}

	// if one of these files exists, probably inside a container
	for _, p := range []string{"/.dockerenv", "/run/.containerenv"} {
		if _, err = os.Stat(p); err == nil {
			return nil
		}
	}

	// Couldn't find anything else, not in a container
	return ErrNotInContainer
}
