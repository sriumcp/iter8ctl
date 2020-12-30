package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
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
	experiment     *v2alpha1.Experiment
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
		experiment:     &v2alpha1.Experiment{},
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

// printAnalysis describes the analysis section of the experiment in a human-interpretable format.
func (d *DescribeCmd) printAnalysis() *DescribeCmd {
	if d.err != nil {
		return d
	}
	sta := d.experiment.Status
	if sta.CompletedIterations == nil {
		fmt.Fprintf(stdout, "Experiment is yet to begin.")
	} else {
		fmt.Fprintf(stdout, "Experiment started. Completed experiment iterations: %v\n", *sta.CompletedIterations)

		if *sta.CompletedIterations > 0 {
			ana := sta.Analysis
			// analysis, _ := json.MarshalIndent(ana, "", "  ")
			// fmt.Fprintln(stdout, string(analysis))

			if ana.WinnerAssessment != nil {
				if ana.WinnerAssessment.Data.WinnerFound {
					fmt.Fprintf(stdout, "Winner: %s\n", *ana.WinnerAssessment.Data.Winner)
				} else {
					fmt.Fprintln(stdout, "No winner found")
				}
			}
		}
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
