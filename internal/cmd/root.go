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
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel = log.InfoLevel
)

var globalFlags = struct {
	LogJSON    bool
	LogLevel   string
	ConfigFile string
}{}

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
		initLogs,
		initConfig,
	)

	rootCmd.PersistentFlags().BoolVar(&globalFlags.LogJSON, "log-json", false, "use JSON formatted logs")
	rootCmd.PersistentFlags().StringVar(&globalFlags.LogLevel, "log-level", defaultLogLevel.String(), "log level")
	rootCmd.PersistentFlags().StringVar(&globalFlags.ConfigFile, "config", "", "config file (default is ~/.config/iscenv/iscenv.yaml")
}

func initLogs() {
	if globalFlags.LogJSON {
		log.SetFormatter(new(log.JSONFormatter))
	} else {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}

	if level, err := log.ParseLevel(globalFlags.LogLevel); err == nil {
		log.SetLevel(level)
		log.WithField("logLevel", level.String()).Debug("Switched log level")
	} else {
		log.WithField("logLevel", globalFlags.LogLevel).Error("Unknown log level, using default")
	}
}

func initConfig() {

	if globalFlags.ConfigFile != "" {
		viper.SetConfigFile(globalFlags.ConfigFile)
	}

	viper.SetConfigName("iscenv.yaml")
	viper.AddConfigPath("$HOME/.config/iscenv/")
	viper.AddConfigPath("./")

	viper.AutomaticEnv()

	// TODO: Error handling here
	viper.ReadInConfig()
}
