package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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
	if d.err != nil {
		return d
	}
	d.err = d.flagSet.Parse(args)
	if d.err != nil {
		d.flagSet.Usage()
		osExiter.Exit(1)
	}
	return d
}

// getExperiment gets the experiment resource object from the k8s cluster
func (d *DescribeCmd) getExperiment() *DescribeCmd {
	if d.err != nil {
		return d
	}
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

func (d *DescribeCmd) printProgress() *DescribeCmd {
	d.description += "******\n"
	d.description += "Experiment name: " + d.experiment.Name + "\n"
	d.description += "Experiment namespace: " + d.experiment.Namespace + "\n"
	d.description += "Experiment target: " + d.experiment.Spec.Target + "\n"
	d.description += "******\n"
	sta := d.experiment.Status
	if sta.CompletedIterations == nil || *sta.CompletedIterations == 0 {
		d.description += "Iteration count: 0\n"
	} else {
		d.description += fmt.Sprintf("Iteration count: %v\n", *sta.CompletedIterations)
	}
	return d
}

func (d *DescribeCmd) printWinnerAssessment() *DescribeCmd {
	if !d.experiment.Started() {
		d.err = errors.New("printWinnerAssessment invoked for experiment that has not started")
		return d
	}
	wa := d.experiment.Status.Analysis.WinnerAssessment
	if wa != nil {
		d.description += "******\n"
		if wa.Data.WinnerFound {
			d.description += fmt.Sprintf("Winner: %s\n", *wa.Data.Winner)
		} else {
			d.description += "Winner: not found\n"
		}
	}
	return d
}

func (d *DescribeCmd) printVersionAssessment() *DescribeCmd {
	if !d.experiment.Started() {
		d.err = errors.New("printVersionAssessment invoked for experiment that has not started")
		return d
	}
	d.description += "******\n"
	versions := d.experiment.GetVersions()
	if d.experiment.RequestCountSpecified() {
		requestCountStrs, err := d.experiment.GetRequestCountStrs()
		if err != nil {
			d.err = err
			return d
		}
		d.description += "Request counts...\n"
		buf := bytes.Buffer{}
		table := tablewriter.NewWriter(&buf)
		table.SetHeader(versions)
		table.Append(requestCountStrs)
		table.Render()
		d.description += buf.String()
	}

	if len(d.experiment.Spec.Criteria.Objectives) > 0 {
		// get criteria list
		// get criteria assessments for versions
		// print
	}
	return d
}

func (d *DescribeCmd) printMetrics() *DescribeCmd {
	sta := d.experiment.Status
	if sta.CompletedIterations == nil || *sta.CompletedIterations == 0 {
		d.description += "Zero experiment iterations complete.\n"
	} else {
		d.description += fmt.Sprintf("%v experiment iterations complete.\n", *sta.CompletedIterations)
	}
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
		// d.printMetrics()
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
