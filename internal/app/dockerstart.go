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
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
)

// Ensure that a container exists and is started.  Returns the ID of the started container or an error
func DockerStart(opts DockerStartOptions) (id string, err error) {
	ilog := log.WithField("instance", opts.Name)
	instances := GetInstances()
	existing := instances.Find(opts.Name)

	// When we recreate we want to maintain the exact same port offset
	if existing != nil {
		ilog.Debug("Found existing instance")
		// Just ensure it's up and return
		if !opts.Recreate {
			ilog.Debug("Starting instance")
			container, err := GetContainerForInstance(existing)
			if err != nil {
				return existing.ID, err
			}
			return existing.ID, DockerClient.StartContainer(existing.ID, container.HostConfig)
		}

		ilog.Debug("Determining existing port offset")
		epo, err := existing.PortOffset()
		if err != nil {
			return "", err
		}

		opts.PortOffset = epo
		opts.PortOffsetSearch = false

		ilog.Debug("Removing instance")
		if err := DockerRemove(existing); err != nil {
			return "", err
		}
		// Reload the instances as the deletion has made the previous list invalid
		instances = GetInstances()
		if err != nil {
			return "", err
		}
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

	ilog.Debug("Creating container")
	container, err := DockerClient.CreateContainer(*opts.ToCreateContainerOptions())
	if err != nil {
		return "", err
	}

	if err := performCopies(container.ID, opts.Copies); err != nil {
		return container.ID, err
	}

	ilog.Debug("Starting container")
	if err = DockerClient.StartContainer(container.ID, opts.ToHostConfig()); err != nil {
		return "", err
	}

	return container.ID, nil
}

func performCopies(id string, copies []string) error {
	r, w := io.Pipe()

	tarErrChan := make(chan error, 1)
	go func() {
		defer w.Close()
		err := writeTar(copies, w)
		log.Debug("Signaling tar waiter channel")
		tarErrChan <- err
	}()

	log.Debug("Uploading copies to container")
	err := DockerClient.UploadToContainer(id, docker.UploadToContainerOptions{
		InputStream:          r,
		Path:                 "/",
		NoOverwriteDirNonDir: true,
	})

	log.Debug("Client complete waiting on tar go routine")
	tarErr := <-tarErrChan

	log.Debug("Tar go routine complete")
	if err != nil || tarErr != nil {
		log.WithError(tarErr).Error("Failed to tar local files for container copy")
		log.WithError(err).Error("Failed to copy local files to container")

		return errors.New("Copy to container failed")
	}

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
