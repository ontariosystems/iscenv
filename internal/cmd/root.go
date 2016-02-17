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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ontariosystems/iscenv/internal/app"
)

var globalFlags = struct {
	Verbose    bool
	ConfigFile string
}{}

var rootCmd = &cobra.Command{
	Use:   "iscenv",
	Short: "Manage Docker-based ISC product environments",
	Long:  "This tool allows the creation and management of Docker-based ISC product Environments.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		app.Fatalf("%s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&globalFlags.Verbose, "verbose", "", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(&globalFlags.ConfigFile, "config", "", "config file (default is ~/.config/iscenv/iscenv.yaml")
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
