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
	Usage          func()
}

// Builder returns an initial Cmd struct pointer.
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

// Error returns any error accumulated by Cmd so far
func (d *Cmd) Error() error {
	return d.err
}

// ParseArgs populates experimentName, experimentNamespace, apiVersion, and kubeconfigPath
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

// GetExperiment gets the experiment resource object from the k8s cluster
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

// printProgress prints the progress of the experiment into the description buffer (d.description)
func (d *Cmd) printProgress() *Cmd {
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

// printProgress prints winner assessment into the description buffer (d.description)
func (d *Cmd) printWinnerAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	wa := d.experiment.Status.Analysis.WinnerAssessment
	if wa != nil {
		d.description.WriteString("\n******\n")
		if wa.Data.WinnerFound {
			d.description.WriteString(fmt.Sprintf("Winning version: %s\n", *wa.Data.Winner))
		} else {
			d.description.WriteString("Winning version: not found\n")
		}
	}
	return d
}

// printObjectiveAssessments is a helper function to print objective assessments for versions into the description buffer (d.description)
func (d *Cmd) printObjectiveAssessments() {
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

// printVersionAssessment prints request counts and criteria assessments into the description buffer (d.description)
func (d *Cmd) printVersionAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	if len(d.experiment.Spec.Criteria.Objectives) > 0 {
		d.printObjectiveAssessments()
	}
	return d
}

func (d *Cmd) printMetrics() *Cmd {
	if d.err != nil {
		return d
	}
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
	return d
}

// PrintAnalysis describes the analysis section of the experiment in a human-interpretable format.
func (d *Cmd) PrintAnalysis() *Cmd {
	if d.err != nil {
		return d
	}
	d.printProgress()
	if d.experiment.Started() {
		d.printWinnerAssessment()
		d.printVersionAssessment()
		d.printMetrics()
	}
	if d.err == nil {
		fmt.Fprintln(d.stdout, d.description.String())
	}
	return d
}
