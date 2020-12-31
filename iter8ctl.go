package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/iter8-tools/iter8ctl/experiment"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

// OSExiter interface enables exiting the current program.
type OSExiter interface {
	Exit(code int)
}
type iter8ctlOS struct{}

func (m iter8ctlOS) Exit(code int) {
	os.Exit(code)
}

var osExiter OSExiter

// Dependency injection for stdio
var stdin io.Reader
var stdout io.Writer

// init initializes stdio, logging, and osExiter
func init() {
	// stdio
	stdin = os.Stdin
	stdout = os.Stdout
	// logging
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
		log.SetLevel(logLevel)
	}
	// osExiter
	osExiter = &iter8ctlOS{}
}

// DescribeCmd struct contains all the data needed for the 'describe' subcommand.
type DescribeCmd struct {
	flagSet        *flag.FlagSet
	experimentPath string
	experiment     *experiment.Experiment
	description    string
	err            error
}

// describeBuilder returns an initial DescribeCmd struct pointer.
func describeBuilder() *DescribeCmd {
	// ContinueOnError enables mocking of os.Exit call in tests with parse errors.
	var flagSet = flag.NewFlagSet("describe", flag.ContinueOnError)

	var d = &DescribeCmd{
		flagSet:        flagSet,
		experimentPath: "",
		experiment:     &experiment.Experiment{},
		description:    "",
		err:            nil,
	}

	//setup flagSet
	const (
		defaultExperimentPath = ""
		usage                 = "absolute path to experiment yaml file, or - for console input (stdin)"
	)
	flagSet.StringVar(&d.experimentPath, "f", defaultExperimentPath, usage)

	return d
}

// parseArgs populates experimentName, experimentNamespace, apiVersion, and kubeconfigPath
func (d *DescribeCmd) parseArgs(args []string) *DescribeCmd {
	d.err = d.flagSet.Parse(args)
	if d.err != nil {
		d.flagSet.Usage()
		osExiter.Exit(1)
	}
	return d
}

// getExperiment gets the experiment resource object from the k8s cluster
func (d *DescribeCmd) getExperiment() *DescribeCmd {
	var expBytes []byte
	if d.experimentPath == "-" {
		expBytes, d.err = ioutil.ReadAll(stdin)
	} else {
		expBytes, d.err = ioutil.ReadFile(d.experimentPath)
	}
	if d.err != nil {
		return d
	}
	expBytesJSON, err := yaml.YAMLToJSON(expBytes)
	d.err = err
	if d.err != nil {
		d.err = errors.New("YAML to JSON conversion error... this could be due to invalid YAML input")
		return d
	}
	d.err = json.Unmarshal(expBytesJSON, d.experiment)
	if d.err != nil {
		d.err = errors.New("unmarshal error... this could be due to invalid experiment YAML input")
		return d
	}
	return d
}

// printProgress prints the progress of the experiment into the description buffer (d.description)
func (d *DescribeCmd) printProgress() *DescribeCmd {
	d.description += "******\n"
	d.description += "Experiment name: " + d.experiment.Name + "\n"
	d.description += "Experiment namespace: " + d.experiment.Namespace + "\n"
	d.description += "Experiment target: " + d.experiment.Spec.Target + "\n"
	d.description += "\n******\n"
	sta := d.experiment.Status
	if sta.CompletedIterations == nil || *sta.CompletedIterations == 0 {
		d.description += "Number of completed iterations: 0\n"
	} else {
		d.description += fmt.Sprintf("Number of completed iterations: %v\n", *sta.CompletedIterations)
	}
	return d
}

// printProgress prints winner assessment into the description buffer (d.description)
func (d *DescribeCmd) printWinnerAssessment() *DescribeCmd {
	wa := d.experiment.Status.Analysis.WinnerAssessment
	if wa != nil {
		d.description += "\n******\n"
		if wa.Data.WinnerFound {
			d.description += fmt.Sprintf("Winning version: %s\n", *wa.Data.Winner)
		} else {
			d.description += "Winning version: not found\n"
		}
	}
	return d
}

// printObjectiveAssessments is a helper function to print objective assessments for versions into the description buffer (d.description)
func (d *DescribeCmd) printObjectiveAssessments() {
	d.description += "\n******\n"
	d.description += "Objectives\n"
	buf := &strings.Builder{}
	table := tablewriter.NewWriter(buf)
	table.SetRowLine(true)
	versions := d.experiment.GetVersions()
	table.SetHeader(append([]string{"Objective"}, versions...))
	for i, objective := range d.experiment.Spec.Criteria.Objectives {
		row := []string{experiment.StringifyObjective(objective)}
		table.Append(append(row, d.experiment.GetSatisfyStrs(i)...))
	}
	table.Render()
	d.description += buf.String()
}

// printVersionAssessment prints request counts and criteria assessments into the description buffer (d.description)
func (d *DescribeCmd) printVersionAssessment() *DescribeCmd {
	if len(d.experiment.Spec.Criteria.Objectives) > 0 {
		d.printObjectiveAssessments()
	}
	return d
}

func (d *DescribeCmd) printMetrics() *DescribeCmd {
	d.description += "\n******\n"
	d.description += "Metrics\n"
	buf := &strings.Builder{}
	table := tablewriter.NewWriter(buf)
	table.SetRowLine(true)
	versions := d.experiment.GetVersions()
	table.SetHeader(append([]string{"Metric"}, versions...))
	for _, metricInfo := range d.experiment.Spec.Metrics {
		row := []string{experiment.GetMetricNameAndUnits(metricInfo)}
		table.Append(append(row, d.experiment.GetMetricValueStrs(metricInfo.Name)...))
	}
	table.Render()
	d.description += buf.String()
	return d
}

// printAnalysis describes the analysis section of the experiment in a human-interpretable format.
func (d *DescribeCmd) printAnalysis() *DescribeCmd {
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
		fmt.Fprintln(stdout, d.description)
	}
	return d
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(stdout, "expected 'describe' subcommand")
		osExiter.Exit(1)
	}

	switch os.Args[1] {

	case "describe":
		d := describeBuilder()
		d.parseArgs(os.Args[2:]).getExperiment().printAnalysis()
		if d.err != nil {
			fmt.Fprintln(stdout, d.err)
			osExiter.Exit(1)
		}

	default:
		fmt.Fprintln(stdout, "expected 'describe' subcommand")
		osExiter.Exit(1)
	}
}
