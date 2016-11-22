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

package cmd

import (
	"github.com/ontariosystems/iscenv/internal/cmd/flags"
	log "github.com/Sirupsen/logrus"
)

func ensureImage() {
	if flags.GetString(rootCmd, "image") == "" {
		logAndExit(log.StandardLogger(), "You must provide an image to use when creating containers either using the --image switch or by setting a value in your configuration file (recommended)")
	}
}
