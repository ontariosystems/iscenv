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
	"strings"
	"time"

	"github.com/ontariosystems/iscenv/v3/internal/cmd/flags"
	"github.com/ontariosystems/iscenv/v3/internal/plugins"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/x-cray/logrus-prefixed-formatter"
)

const (
	defaultLogLevel           = log.InfoLevel
	defaultConfigFile         = "iscenv"
	defaultConfigDir          = "$HOME/.config/iscenv/"
	disablePrimaryCommandFile = "/.iscenv-disable-primary-command"
)

var rootCmd = &cobra.Command{
	Use:   "iscenv",
	Short: "Manage Docker-based ISC product environments",
	Long:  "This tool allows the creation and management of Docker-based ISC product Environments.",
}

// Execute runs the root functionality of the iscenv command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logAndExit(log.WithError(err), "Failed to execute iscenv command")
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

	// The following flags are not used by every command but are used by multiple commands and do not make sense to differ between those commands (in fact, it would be detrimental).

	// This flag is the image which will be used by default when creating new containers or listing versions.
	flags.AddFlagComplex(rootCmd, true, true, "image", "", "", "the image to use when creating ISC product containers.  You will want to set a default for this in your configuration file (eg. mycompany/ensemble)")

	// This allows us to set lifecycle plugins across the multiple commands to which they belong will still allowing versioners to use their own set of plugins
	if !skipPluginActivation() {
		// Logging can't have been configured yet, so we're using an empty PluginArgs
		var lcs []*plugins.ActivatedLifecycler
		defer getActivatedLifecyclers(nil, plugins.PluginArgs{}, &lcs)()
		available := make([]string, len(lcs))
		for i, lc := range lcs {
			available[i] = lc.Id
		}
		flags.AddFlagComplex(rootCmd, true, true, "plugins", "", "", "An ordered comma-separated list of lifecycle plugins you wish to activate. available plugins: "+strings.Join(available, ","))
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
		if !os.IsNotExist(err) && !configFileNotFound(err) {
			fmt.Printf("Invalid configuration file, err: %s\n", err)
			os.Exit(1)
		}
	}
}

func configFileNotFound(err error) bool {
	_, ok := err.(viper.ConfigFileNotFoundError)
	return ok
}

func initLogs() {
	// Unfortunately, we cannot log in this function or it breaks plugins by sending non-plugin api messages across the pipe
	// The default formatter is JSON
	if !flags.GetBool(rootCmd, "log-json") {
		log.SetFormatter(&prefixed.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.Stamp,
		})
	}

	levelString := flags.GetString(rootCmd, "log-level")
	if level, err := log.ParseLevel(levelString); err == nil {
		log.SetLevel(level)
	}
}
