/*
Copyright 2015 Ontario Systems

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

package main

import (
	"io/ioutil"
	"strings"
)

func ensureWithinContainer(commandName string) {

	proc1CGroupContents, err := ioutil.ReadFile("/proc/1/cgroup")

	if err != nil {
		fatalf("Could not read /proc/1/cgroup to determine environment")
	}

	// if we have some control groups owned by docker, then we are within a container
	contents := string(proc1CGroupContents)
	if !strings.Contains(contents, ":/docker/") {
		fatalf("Cannot run `%s` outside of a container!\n", commandName)
	}
}
