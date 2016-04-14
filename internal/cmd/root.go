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
	"fmt"
	"os"

	"github.com/ontariosystems/iscenv/internal/cmd/flags"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel   = log.InfoLevel
	defaultConfigFile = "iscenv"
	defaultConfigDir  = "$HOME/.config/iscenv/"
)

var rootCmd = &cobra.Command{
	Use:   "iscenv",
	Short: "Manage Docker-based ISC product environments",
	Long:  "This tool allows the creation and management of Docker-based ISC product Environments.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initLogs,
	)

	flags.AddFlagComplex(rootCmd, true, false, "config", "", "", "the path to a configuration file in json, toml, yaml or hcl format")
	flags.AddFlagComplex(rootCmd, true, true, "log-level", "", defaultLogLevel.String(), "log level")
	flags.AddFlagComplex(rootCmd, true, true, "log-json", "", false, "use JSON formatted logs")

	// This flag is the image which will be used by default when creating new containers or listing versions.
	// Technically, very few commands actually need this flag but I did not want the config to have to have it specified multiple places when it doesn't make sense for it to differ between commands (and would actually be detrimental)
	flags.AddFlagComplex(rootCmd, true, true, "image", "", "", "the image to use when creating ISC product containers.  You will want to set a default for this in your configuration file (eg. mycompany/ensemble)")
}

func initLogs() {
	if flags.GetBool(rootCmd, "log-json") {
		log.SetFormatter(new(log.JSONFormatter))
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}

	levelString := flags.GetString(rootCmd, "log-level")
	if level, err := log.ParseLevel(levelString); err == nil {
		log.SetLevel(level)
		log.WithField("logLevel", level.String()).Debug("Switched log level")
	} else {
		log.WithField("logLevel", levelString).Error("Unknown log level, using default")
	}
}

func initConfig() {
	if configPath := flags.GetString(rootCmd, "config"); configPath != "" {
		viper.SetConfigFile(configPath)
	}

	viper.SetConfigName(defaultConfigFile)
	viper.AddConfigPath(defaultConfigDir)
	viper.AddConfigPath("./")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) && !unsupportedConfigTypeBlank(err) {
			fmt.Printf("Invalid configuration file, err: %s\n", err)
			os.Exit(1)
		}
	}
}

// TODO: Remove this code when https://github.com/spf13/viper/pull/161 is merged and the updated version is vendored
func unsupportedConfigTypeBlank(err error) bool {
	return err.Error() == `Unsupported Config Type ""`
}