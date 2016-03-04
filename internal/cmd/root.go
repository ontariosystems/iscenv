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
	defaultConfigExt  = "json"
	defaultConfigDir  = "$HOME/.config/iscenv/"
	defaultConfigPath = defaultConfigDir + defaultConfigFile + "." + defaultConfigExt
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

	flags.AddFlagComplex(rootCmd, true, false, "config", "", "", "config file (default is "+defaultConfigPath+")")
	flags.AddFlagComplex(rootCmd, true, true, "log-level", "", defaultLogLevel.String(), "log level")
	flags.AddFlagComplex(rootCmd, true, true, "log-json", "", false, "use JSON formatted logs")
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

	viper.SetConfigType(defaultConfigExt)
	viper.SetConfigName(defaultConfigFile)
	viper.AddConfigPath(defaultConfigDir)
	viper.AddConfigPath("./")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Invalid configuration file, err: %s\n", err)
			os.Exit(1)
		}
	}
}
