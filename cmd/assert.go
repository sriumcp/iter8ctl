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

	"github.com/spf13/cobra"
)

var conditions []string

// assertCmd represents the assert command
var assertCmd = &cobra.Command{
	Use:   "assert -c experiment-condition ...",
	Short: "Assert conditions for the experiment",
	Long:  `One or more conditions can be asserted using this command for an Iter8 experiment. This command is especially in CI/CD/Gitops pipelines to assert conditions prior to version promotion or rollback.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assert called")
	},
}

func init() {
	rootCmd.AddCommand(assertCmd)
	assertCmd.Flags().StringSliceVarP(&conditions, "condition", "c", nil, "completed | failure | handlerFailure | successful | candidateWon | baselineWon | noWinner | winnerFound")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
