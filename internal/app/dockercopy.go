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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

func DockerCopy(instanceName, instancePath, localPath string) error {
	instance := GetInstances().Find(strings.ToLower(instanceName))
	if instance == nil {
		return fmt.Errorf("Could not find instance, name: %s", instanceName)
	}

	r, w := io.Pipe()

	go func() {
		DockerClient.DownloadFromContainer(instance.ID, docker.DownloadFromContainerOptions{
			Path:         instancePath,
			OutputStream: w,
		})
	}()

	t := tar.NewReader(r)
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(localPath, header.Name)
		info := header.FileInfo()
		fmt.Println(path, info.Name())
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, t)
		if err != nil {
			return err
		}
	}

	return nil
}
