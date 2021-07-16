/* Package cmd with describe subcommand */
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
