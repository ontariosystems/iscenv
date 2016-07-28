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

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/pflag"
)

func ServeLifecyclePlugin(impl Lifecycler) {
	pluginMap := map[string]plugin.Plugin{
		LifecyclerKey: LifecyclerPlugin{Plugin: impl},
	}

	// See configureLogger comments
	go configureLogger(os.Stdout)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: PluginHandshake,
		Plugins:         pluginMap,
	})
}

func ServeVersionsPlugin(impl Versioner) {
	pluginMap := map[string]plugin.Plugin{
		VersionerKey: VersionerPlugin{Plugin: impl},
	}

	// See configureLogger comments
	go configureLogger(os.Stdout)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: PluginHandshake,
		Plugins:         pluginMap,
	})
}

// plugin.Serve switches stdout/stderr and the logger will not work unless configured after this.
// plugin.Serve also blocks so we can't run the configuration afterwards.
// The options which presented themselves were...
//  - Configure the logging in every "event" handler
//  - Configure the logging once by polling for a stdout change and then triggering the configuration.  I chose this one.  However, it's brittle and if the logging breaks you should look here
func configureLogger(oldOut io.Writer) {
	for {
		if oldOut != os.Stdout {
			break
		}
	}

	log.SetOutput(os.Stdout)

	var l *string = pflag.String("log-level", "", "")
	var j *bool = pflag.Bool("log-json", false, "")
	pflag.Parse()

	if *j {
		log.SetFormatter(new(log.JSONFormatter))
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}

	if level, err := log.ParseLevel(*l); err == nil {
		log.SetLevel(level)
	}
}
