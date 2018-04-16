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
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	version "github.com/hashicorp/go-version"
)

// Ensure that a container exists and is started.  Returns the ID of the started container or an error
func DockerStart(opts DockerStartOptions) (id string, err error) {
	ilog := log.WithField("instance", opts.Name)
	instances := GetInstances()
	existing := instances.Find(opts.Name)

	var container *docker.Container
	if existing != nil {
		ilog.Debug("Found existing instance")
		ilog.Debug("Determining existing port offset")
		opts.PortOffsetSearch = false
		opts.PortOffset, err = existing.PortOffset()
		if err != nil {
			return "", err
		}

		if opts.Recreate {
			ilog.Info("Removing existing instance")
			if err := DockerRemove(existing); err != nil {
				return "", err
			}

			// Reload the instances as the deletion has made the previous list invalid
			instances = GetInstances()
			if err != nil {
				return "", err
			}

			existing = nil
		} else if container, err = GetContainerForInstance(existing); err != nil {
			return "", err
		}
	}

	if existing == nil {
		if opts.PortOffsetSearch {
			if opts.PortOffset, err = instances.CalculatePortOffset(opts.PortOffset); err != nil {
				return "", err
			}
		} else {
			if upo, err := instances.UsedPortOffset(opts.PortOffset); err != nil {
				return "", err
			} else if upo && !opts.DisablePortOffsetConflictCheck {
				return "", fmt.Errorf("Port offset conflict, offset: %d", opts.PortOffset)
			}
		}

		containerOpts := *opts.ToCreateContainerOptions()
		b, err := json.Marshal(containerOpts)
		if err != nil {
			return "", err
		}

		ilog.WithField("opts", string(b)).Debug("Creating container")
		container, err = DockerClient.CreateContainer(containerOpts)
		if err != nil {
			return "", err
		}
	}

	if container.State.Running {
		ilog.Info("Instance is already running")
		return existing.ID, nil
	}

	if err := performCopies(container.ID, opts.Copies); err != nil {
		return container.ID, err
	}

	ilog.Debug("Starting container")
	hostConfig, err := getStartHostConfig(opts)
	if err != nil {
		return "", err
	}
	if err = DockerClient.StartContainer(container.ID, hostConfig); err != nil {
		return "", err
	}

	return container.ID, nil
}

// Passing HostConfig is deprecated after 1.10 and removed at 1.12 but prior to that we had problems if we didn't pass it.
// As such, we will return a nil on newer versions but the actual hostconfig on older ones.
func getStartHostConfig(opts DockerStartOptions) (*docker.HostConfig, error) {
	cutoff, _ := version.NewVersion("1.10.0")

	env, err := DockerClient.Version()
	if err != nil {
		return nil, err
	}

	ver, err := version.NewVersion(env.Get("Version"))
	if err != nil {
		return nil, err
	}

	if ver.LessThan(cutoff) {
		return opts.ToHostConfig(), nil
	}

	return nil, nil
}

func performCopies(id string, copies []string) error {
	r, w := io.Pipe()

	tarErrChan := make(chan error, 1)
	go func() {
		defer w.Close()
		tarErrChan <- writeTar(copies, w)
	}()

	log.Debug("Uploading copies to container")
	err := DockerClient.UploadToContainer(id, docker.UploadToContainerOptions{
		InputStream:          r,
		Path:                 "/",
		NoOverwriteDirNonDir: true,
	})

	if err != nil {
		log.WithError(err).Error("Failed to copy local files to container")
		return err
	}

	log.Debug("Client complete waiting on tar go routine")
	if err := <-tarErrChan; err != nil {
		log.WithError(err).Error("Failed to tar local files for container copy")
		return err
	}

	log.Debug("Tar go routine complete")

	return nil
}

func writeTar(copies []string, writer io.Writer) error {
	tw := tar.NewWriter(writer)

	for _, copy := range copies {
		s := strings.Split(copy, ":")
		source := s[0]
		dest := s[1]
		log.WithFields(log.Fields{"source": source, "destination": dest}).Debugf("Processing copy")
		sourceFile, err := os.Open(source)
		if err != nil {
			return err
		}

		fi, err := sourceFile.Stat()
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: dest,
			Mode: int64(fi.Mode()),
			Size: fi.Size(),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := io.Copy(tw, sourceFile); err != nil {
			return err
		}

		sourceFile.Close()
	}

	log.Debug("Closing tar writer")
	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}
