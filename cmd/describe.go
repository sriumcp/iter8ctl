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
	"errors"

	"github.com/iter8-tools/iter8ctl/describe"
	expr "github.com/iter8-tools/iter8ctl/experiment"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use: "describe",
	Args: func(cmd *cobra.Command, args []string) error {
		if experiment == "" && !latest {
			return errors.New("Either specify a valid experiment name with -e or use the latest option with -l")
		}
		var err error
		if exp, err = expr.GetExperiment(latest, experiment, namespace); err != nil {
			return err
		}
		return nil
	},
	Short: "Describe an Iter8 experiment",
	Long:  `Summarize an experiment, including the stage of the experiment, how versions are performing with respect to the experiment criteria (reward, SLOs, metrics), and information about the winning version. This program is a K8s client and requires a valid K8s cluster with Iter8 installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		describe.Builder().WithExperiment(exp).PrintAnalysis()
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// describeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// describeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}