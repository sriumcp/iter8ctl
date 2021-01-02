// Package describe implements the `iter8ctl describe` subcommand.
package describe

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/iter8-tools/iter8ctl/experiment"
	"github.com/olekukonko/tablewriter"
	"sigs.k8s.io/yaml"
)

// Cmd struct contains all the data needed for the 'describe' subcommand.
type Cmd struct {
	flagSet        *flag.FlagSet
	experimentPath string
	experiment     *experiment.Experiment
	description    strings.Builder
	err            error
	stdin          io.Reader
	stdout         io.Writer
	stderr         io.Writer
	// Usage is a function that is invoked when execution of any `Cmd` method results in an error.
	// Typically, Usage() prints the error message to stderr, and the program exits.
	// Note that Usage() is a struct field and not a method.
	// You can supply your own implementation of Usage() while constructing a new `Cmd` struct.
	Usage func()
}

// Builder returns an initialized Cmd struct pointer.
func Builder(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Cmd {
	// ContinueOnError ensures flagSet does not ExitOnError
	var flagSet = flag.NewFlagSet("describe", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	var d = &Cmd{
		flagSet:        flagSet,
		experimentPath: "",
		experiment:     &experiment.Experiment{},
		description:    strings.Builder{},
		err:            nil,
		stdin:          stdin,
		stdout:         stdout,
		stderr:         stderr,
		Usage:          flagSet.Usage,
	}

	//setup flagSet
	const (
		defaultExperimentPath = ""
		usage                 = "absolute path to experiment yaml file, or - for console input (stdin)"
	)
	flagSet.StringVar(&d.experimentPath, "f", defaultExperimentPath, usage)

	return d
}

// Error returns any error accumulated by Cmd so far.
func (d *Cmd) Error() error {
	return d.err
}

// ParseArgs populates d.experimentPath.
func (d *Cmd) ParseArgs(args []string) *Cmd {
	if d.err != nil {
		return d
	}
	d.err = d.flagSet.Parse(args)
	if d.err == nil {
		if d.flagSet.NFlag() == 0 {
			d.err = errors.New("Missing experiment")
			fmt.Fprintln(d.stderr, d.err)
			d.Usage()
		}
	}
	return d
}

// GetExperiment populates d.experiment from an input file or from stdin input.
// Input should be valid experiment YAML.
// If input is invalid, GetExperiment sets an error in d.err.
func (d *Cmd) GetExperiment() *Cmd {
	if d.err != nil {
		return d
	}
	var expBytes []byte
	if d.experimentPath == "-" {
		expBytes, d.err = ioutil.ReadAll(d.stdin)
	} else {
		expBytes, d.err = ioutil.ReadFile(d.experimentPath)
	}
	if d.err != nil {
		d.err = errors.New("Error reading experiment YAML input")
		fmt.Fprintln(d.stderr, d.err)
		return d
	}
	expBytesJSON, err := yaml.YAMLToJSON(expBytes)
	d.err = err
	if d.err != nil {
		d.err = errors.New("YAML to JSON conversion error... this could be due to invalid YAML input")
		fmt.Fprintln(d.stderr, d.err)
		return d
	}
	d.err = json.Unmarshal(expBytesJSON, d.experiment)
	if d.err != nil {
		d.err = errors.New("unmarshal error... this could be due to invalid experiment YAML input")
		fmt.Fprintln(d.stderr, d.err)
		return d
	}
	return d
}

// PrintProgress prints name, namespace, and target of the experiment and the number of completed iterations.
func (d *Cmd) PrintProgress() *Cmd {
	if d.err != nil {
		return d
	}
	d.description.WriteString("******\n")
	d.description.WriteString("Experiment name: " + d.experiment.Name + "\n")
	d.description.WriteString("Experiment namespace: " + d.experiment.Namespace + "\n")
	d.description.WriteString("Experiment target: " + d.experiment.Spec.Target + "\n")
	d.description.WriteString("\n******\n")
	sta := d.experiment.Status
	if sta.CompletedIterations == nil || *sta.CompletedIterations == 0 {
		d.description.WriteString("Number of completed iterations: 0\n")
	} else {
		d.description.WriteString(fmt.Sprintf("Number of completed iterations: %v\n", *sta.CompletedIterations))
	}
	return d
}

// PrintWinnerAssessment prints the winning version in the experiment, if Status.Analysis.WinnerAssessment is not nil in the experiment object.
func (d *Cmd) PrintWinnerAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	if a := d.experiment.Status.Analysis; a != nil {
		if w := a.WinnerAssessment; w != nil {
			d.description.WriteString("\n******\n")
			if w.Data.WinnerFound {
				d.description.WriteString(fmt.Sprintf("Winning version: %s\n", *w.Data.Winner))
			} else {
				d.description.WriteString("Winning version: not found\n")
			}
		}
	}
	return d
}

// PrintObjectiveAssessment prints a matrix of boolean values, if Status.Analysis.VersionAssessments is not nil in the experiment object.
// Rows correspond to experiment objectives, columns correspond to versions, and entry [i, j] indicates if version j satisfies objective i.
// Objective assessments are printed in the same sequence as in the experiment's spec.criteria.objectives section.
func (d *Cmd) PrintObjectiveAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	if a := d.experiment.Status.Analysis; a != nil {
		if v := a.VersionAssessments; v != nil {
			d.description.WriteString("\n******\n")
			d.description.WriteString("Objectives\n")
			table := tablewriter.NewWriter(&d.description)
			table.SetRowLine(true)
			versions := d.experiment.GetVersions()
			table.SetHeader(append([]string{"Objective"}, versions...))
			for i, objective := range d.experiment.Spec.Criteria.Objectives {
				row := []string{experiment.StringifyObjective(objective)}
				table.Append(append(row, d.experiment.GetSatisfyStrs(i)...))
			}
			table.Render()
		}
	}
	return d
}

// PrintVersionAssessment prints how each version is performing with respect to experiment criteria.
func (d *Cmd) PrintVersionAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	if c := d.experiment.Spec.Criteria; c != nil && len(c.Objectives) > 0 {
		d.PrintObjectiveAssessment()
	}
	return d
}

// PrintMetrics prints a matrix of decimal values.
// Rows correspond to experiment metrics, columns correspond to versions, and entry [i, j] indicates the value of metric i for version j.
// Metrics are in the same sequence as in the experiment's spec.metrics section.
func (d *Cmd) PrintMetrics() *Cmd {
	if d.err != nil {
		return d
	}
	if a := d.experiment.Status.Analysis; a != nil {
		if v := a.AggregatedMetrics; v != nil {
			d.description.WriteString("\n******\n")
			d.description.WriteString("Metrics\n")
			table := tablewriter.NewWriter(&d.description)
			table.SetRowLine(true)
			versions := d.experiment.GetVersions()
			table.SetHeader(append([]string{"Metric"}, versions...))
			for _, metricInfo := range d.experiment.Spec.Metrics {
				row := []string{experiment.GetMetricNameAndUnits(metricInfo)}
				table.Append(append(row, d.experiment.GetMetricValueStrs(metricInfo.Name)...))
			}
			table.Render()
		}
	}
	return d
}

// PrintAnalysis describes the experiment by printing progress, winner and version assessments, and metrics.
func (d *Cmd) PrintAnalysis() *Cmd {
	if d.err != nil {
		return d
	}
	d.PrintProgress()
	if d.experiment.Started() {
		d.PrintWinnerAssessment().PrintVersionAssessment().PrintMetrics()
	}
	if d.err == nil {
		fmt.Fprintln(d.stdout, d.description.String())
	}
	return d
}
