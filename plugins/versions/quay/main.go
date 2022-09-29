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

package quayversionsplugin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/v3/iscenv"
	log "github.com/sirupsen/logrus"
)

var plog = log.WithField("plugin", pluginKey)

const (
	pluginKey   = "quay"
	registry    = "quay.io"
	registryURL = "https://" + registry
)

// Plugin represents this plugin and serves as a place to attach functions to implement the Lifecycler interface
type Plugin struct{}

// Main serves as the main entry point for the plugin
func (plugin *Plugin) Main() {
	iscenv.ServeVersionsPlugin(plugin)
}

// Key returns the unique identifier for the plugin
func (*Plugin) Key() string {
	return pluginKey
}

// Versions finds the versions available for the provided image
func (*Plugin) Versions(image string) (iscenv.ISCVersions, error) {
	// If this isn't a quay image, we can't search quay for it (obviously)
	prefix := registry + "/"
	if !strings.HasPrefix(image, prefix) {
		return nil, nil
	}

	img := strings.TrimPrefix(image, prefix)

	authcfg, err := docker.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		return nil, err
	}

	auth := authcfg.Configs["quay.io"]
	ru, _ := url.Parse(registryURL)
	client := newRegistryClient(ru, auth.Username, auth.Password)
	// This is the URL for a v2 repo.  Quay does not appear to have this route implemented
	//reqURL := fmt.Sprintf("%s/v2/%s/tags/list", registryURL, img)
	// So, we are using the deprecated 1.0 route. This makes the v2 auth unnecessary but
	// it still works as expected and leaves us ready to quickly switch to a v2 route.
	reqURL := fmt.Sprintf("%s/v1/repositories/%s/tags", registryURL, img)
	plog.WithField("url", reqURL).Debug("Requesting tag list from API")
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	// The json is just a map of version to image id
	// Unfortunately, the image ID is wrong (thanks quay), so, we're just gonna say "remote"
	images := make(map[string]string)
	if err := dec.Decode(&images); err != nil {
		return nil, err
	}

	versions := make(iscenv.ISCVersions, len(images))
	i := 0
	for v, id := range images {
		versions[i] = &iscenv.ISCVersion{
			ID:      id,
			Version: v,
			Created: time.Now().Unix(),
			Source:  "quay.io",
		}
		i++
	}

	return versions, nil
}
