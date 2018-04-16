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
	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/iscenv/iscenv"
)

func InstanceLogger(instance *iscenv.ISCInstance) *log.Entry {
	return InstanceLoggerArgs(instance.Name, instance.ID)
}

// Return an instance logger with the values as args
func InstanceLoggerArgs(instanceName, instanceID string) *log.Entry {
	return log.WithFields(log.Fields{
		"instanceName": instanceName,
		"instanceID":   instanceID,
	})
}

func PluginLogger(id, method, path string) *log.Entry {
	return log.WithFields(log.Fields{
		"pluginID":   id,
		"method":     method,
		"pluginPath": path,
	})
}

func DockerRepoLogger(repo string) *log.Entry {
	return log.WithField("dockerRepository", repo)
}

// Will evaluate an error for known error types and return a logger with the appropriate fields pulled from that type of error.
func ErrorLogger(logger log.FieldLogger, err error) *log.Entry {
	switch e := err.(type) {
	case *PluginError:
		logger = logger.WithFields(log.Fields{
			"pluginID":   e.Plugin,
			"method":     e.PluginMethod,
			"pluginPath": e.PluginPath,
		})
		err = e.Err
	case *InstanceError:
		logger = logger.WithFields(log.Fields{
			"instanceName": e.InstanceName,
			"instanceID":   e.InstanceID,
		})
		err = e.Err
	case *DockerConfigError:
		logger = logger.WithField("configPath", e.ConfigPath)
		if e.Registry != "" {
			logger = logger.WithField("registry", e.Registry)
		}
		err = e.Err
	}

	return logger.WithError(err)
}
