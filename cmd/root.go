// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pgm",
	Short: "Pokemon Go Mapper",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().String("config", "", "config file (default is pgm.yaml)")

	RootCmd.PersistentFlags().String("maps", "", "Google maps API Key")
	viper.BindPFlag("maps", RootCmd.PersistentFlags().Lookup("maps"))

	RootCmd.PersistentFlags().String("location", "", "Starting location")
	viper.BindPFlag("location", RootCmd.PersistentFlags().Lookup("location"))

	RootCmd.PersistentFlags().Int("steps", 5, "The step radius of the beehive hexagon (radiusMeters = steps * 70)")
	viper.BindPFlag("steps", RootCmd.PersistentFlags().Lookup("steps"))

	RootCmd.PersistentFlags().Int("jitter", 9, "Jitter the starting location by up to this many meters in a random direction")
	viper.BindPFlag("jitter", RootCmd.PersistentFlags().Lookup("jitter"))

	RootCmd.PersistentFlags().StringSlice("accounts", nil, "accounts in the form of ptc:user:pass or google:user:pass")
	viper.BindPFlag("accounts", RootCmd.PersistentFlags().Lookup("accounts"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetEnvPrefix("pgm")
	viper.SetConfigName("pgm")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("/etc/")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
