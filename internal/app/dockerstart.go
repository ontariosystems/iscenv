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
)

// Ensure that a container exists and is started.  Returns the ID of the started container or an error
func DockerStart(opts DockerStartOptions) (id string, err error) {
	instances := GetInstances()
	existing := instances.Find(opts.Name)

	// When we recreate we want to maintain the exact same port offset
	if existing != nil {
		// Just ensure it's up and return
		if !opts.Recreate {
			err := DockerClient.StartContainer(existing.ID, GetContainerForInstance(existing).HostConfig)
			return existing.ID, err
		}

		epo, err := existing.PortOffset()
		if err != nil {
			return "", err
		}

		opts.PortOffset = epo
		opts.PortOffsetSearch = false

		if _, err := DockerRemove(existing.Name); err != nil {
			return "", err
		}
		// Reload the instances as the deletion has made the previous list invalid
		instances = GetInstances()
	}

	if opts.PortOffsetSearch {
		if opts.PortOffset, err = instances.CalculatePortOffset(opts.PortOffset); err != nil {
			return "", err
		}
	} else {
		if upo, err := instances.UsedPortOffset(opts.PortOffset); err != nil {
			return "", err
		} else if upo {
			return "", fmt.Errorf("Port offset conflict, offset: %s", opts.PortOffset)
		}
	}

	x := opts.ToCreateContainerOptions()
	fmt.Printf("LEH: %s\n", x.Config.Image)

	container, err := DockerClient.CreateContainer(*opts.ToCreateContainerOptions())
	if err != nil {
		return "", err
	}

	if err = DockerClient.StartContainer(container.ID, opts.ToHostConfig()); err != nil {
		return "", err
	}

	return container.ID, nil
}
