package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/iter8-tools/iter8ctl/describe"
)

// OSExiter interface enables exiting the current program.
// The inferface is useful in tests to mock OS.exit function with GO's panic function.
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
var stderr io.Writer

// init initializes stdio, logging, and osExiter
func init() {
	// stdio
	stdin = os.Stdin
	stdout = os.Stdout
	stderr = os.Stderr
	// logging
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(stdout)
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
		log.SetLevel(logLevel)
	}
	// osExiter
	osExiter = &iter8ctlOS{}
}

// main serves as the program entry point.
func main() {
	d := describe.Builder(stdin, stdout, stderr)
	if len(os.Args) < 2 {
		fmt.Fprintln(stderr, "expected 'describe' subcommand")
		d.Usage()
		osExiter.Exit(1)
	}

	switch os.Args[1] {

	case "describe":
		d.ParseArgs(os.Args[2:]).GetExperiment().PrintAnalysis()
		if d.Error() != nil {
			osExiter.Exit(1)
		}

	default:
		fmt.Fprintln(stderr, "expected 'describe' subcommand")
		d.Usage()
		osExiter.Exit(1)
	}
}
