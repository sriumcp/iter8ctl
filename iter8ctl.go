package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

// init initializes logging, osExiter, and flagSet
func init() {
	// logging
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
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
	experimentPath *string
	experiment     *v2alpha1.Experiment
	err            error
}

// describeBuilder returns an initial DescribeCmd struct pointer.
func describeBuilder() *DescribeCmd {
	var flagSet = flag.NewFlagSet("describe", flag.ContinueOnError)

	var experimentPath string

	//setup flagSet
	const (
		defaultExperimentPath = "-"
		usage                 = "absolute path to experiment yaml filename or - for console input (stdin)"
	)
	flagSet.StringVar(&experimentPath, "e", defaultExperimentPath, usage)

	return &DescribeCmd{
		flagSet:        flagSet,
		experimentPath: &experimentPath,
		experiment:     &v2alpha1.Experiment{},
		err:            nil,
	}
}

// parseArgs populates experimentName, experimentNamespace, apiVersion, and kubeconfigPath
func (d *DescribeCmd) parseArgs(args []string) *DescribeCmd {
	if d.err != nil {
		return d
	}
	d.err = d.flagSet.Parse(args)
	if d.err != nil {
		d.flagSet.Usage()
	}
	return d
}

// getExperiment gets the experiment resource object from the k8s cluster
func (d *DescribeCmd) getExperiment() *DescribeCmd {
	if d.err != nil {
		return d
	}
	var expBytes []byte
	if *d.experimentPath == "-" {
		expBytes, d.err = ioutil.ReadAll(os.Stdin)
	} else {
		expBytes, d.err = ioutil.ReadFile(*d.experimentPath)
	}
	if d.err != nil {
		return d
	}
	expBytesJSON, err := yaml.YAMLToJSON(expBytes)
	d.err = err
	if d.err != nil {
		return d
	}
	d.err = json.Unmarshal(expBytesJSON, d.experiment)
	if d.err != nil {
		return d
	}
	return d
}

// printAnalysis describes the analysis section of the experiment in a human-interpretable format.
func (d *DescribeCmd) printAnalysis() *DescribeCmd {
	if d.err != nil {
		return d
	}
	return d
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'describe' subcommand")
		osExiter.Exit(1)
	}

	switch os.Args[1] {

	case "describe":
		d := describeBuilder()
		d.parseArgs(os.Args[2:]).getExperiment().printAnalysis()
		if d.err != nil {
			osExiter.Exit(1)
		}

	default:
		fmt.Println("expected 'describe' subcommand")
		osExiter.Exit(1)
	}
}
