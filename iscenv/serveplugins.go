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
	"io"
	"os"
	"time"

	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/x-cray/logrus-prefixed-formatter"
)

// ServeLifecyclePlugin serves a life cycle plugin
func ServeLifecyclePlugin(impl Lifecycler) {
	serve(map[string]plugin.Plugin{
		LifecyclerKey: LifecyclerPlugin{Plugin: impl},
	})
}

// ServeVersionsPlugin serves a versioner plugin
func ServeVersionsPlugin(impl Versioner) {
	serve(map[string]plugin.Plugin{
		VersionerKey: VersionerPlugin{Plugin: impl},
	})
}

func serve(pluginMap map[string]plugin.Plugin) {
	configureLogger()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: PluginHandshake,
		Plugins:         pluginMap,
	})
}

// plugin.Serve switches stdout/stderr and the logger will not work unless configured after this.
// plugin.Serve also blocks so we can't run the configuration afterwards.
// The options which presented themselves were...
//   - Configure the logging in every "event" handler
//   - Configure the logging once by polling for a stdout change and then triggering the configuration.  I chose this one.  However, it's brittle and if the logging breaks you should look here
func configureLogger() {
	// discard logs until the correct device exists
	log.SetOutput(io.Discard)

	// go func to watch for os.Stdout to change
	oldOut := os.Stdout
	go func() {
		for {
			if oldOut != os.Stdout {
				break
			}
		}

		log.SetOutput(os.Stdout)
	}()

	l := pflag.String("log-level", "info", "")
	j := pflag.Bool("log-json", false, "")
	pflag.Parse()

	if *j {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339})
	} else {
		log.SetFormatter(&prefixed.TextFormatter{ForceColors: true})
	}

	if level, err := log.ParseLevel(*l); err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.InfoLevel)
		log.WithError(err).Error("Could not set log level")
	}
}
