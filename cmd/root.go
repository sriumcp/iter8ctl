/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	expr "github.com/iter8-tools/iter8ctl/experiment"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var experiment string
var namespace string
var latest bool
var exp *expr.Experiment

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iter8ctl",
	Short: "Iter8 command line utility",
	Long:  `iter8ctl promotes understanding of an Iter8 experiment. It can be used to describe the stage of the experiment, how versions are performing, and assert various conditions relating to the experiment. This program is a K8s client and requires a valid K8s cluster with Iter8 installed.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iter8ctl.yaml)")

	rootCmd.PersistentFlags().StringVarP(&experiment, "experiment", "e", "", "name of the experiment; ignored when -l flag is used")

	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "namespace of the experiment; ignored when -l flag is used")

	rootCmd.PersistentFlags().BoolVarP(&latest, "latest", "l", false, "use the Iter8 experiment with the latest creation timestamp; either specify this flag or use -e")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".iter8ctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".iter8ctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}