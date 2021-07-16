/*
Copyright © 2021 Iter8

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
	"fmt"
	"os"

	expr "github.com/iter8-tools/iter8ctl/experiment"
	"github.com/spf13/cobra"
)

var conditions []string
var conds []expr.ConditionType

// assertCmd represents the assert command
var assertCmd = &cobra.Command{
	Use: "assert",
	Args: func(cmd *cobra.Command, args []string) error {
		if experiment == "" && !latest {
			return errors.New("Either specify a valid experiment name with -e or use the latest option with -l")
		}
		if conditions == nil || len(conditions) == 0 {
			return errors.New("One or more conditions must be specified with assert")
		}
		for _, cond := range conditions {
			switch cond {
			case string(expr.Completed):
				conds = append(conds, expr.Completed)
			case string(expr.WinnerFound):
				conds = append(conds, expr.WinnerFound)
			default:
				return errors.New("Invalid condition: " + cond)
			}
		}
		return nil
	},
	Short: "Assert conditions for the experiment",
	Long:  `One or more conditions can be asserted using this command for an Iter8 experiment. This command is especially useful in CI/CD/Gitops pipelines prior to version promotion or rollback. This program is a K8s client and requires a valid K8s cluster with Iter8 installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := exp.Assert(conds); err == nil {
			fmt.Println("All conditions satisfied.")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(assertCmd)
	assertCmd.Flags().StringSliceVarP(&conditions, "condition", "c", nil, "completed | winnerFound")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
