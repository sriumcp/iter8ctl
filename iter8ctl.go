package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/iter8-tools/iter8ctl/describe"
)

// stdin enables dependency injection for console input (stdin)
var stdin io.Reader

// stdout enables dependency injection for console output (stdout)
var stdout io.Writer

// stderr enables dependency injection for console error output (stderr)
var stderr io.Writer

// init initializes stdin/out/err, and logging.
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
}

// main serves as the entry point for iter8ctl CLI.
func main() {
	d := describe.Builder(stdin, stdout, stderr)
	if len(os.Args) < 2 {
		fmt.Fprintln(stderr, "expected 'describe' subcommand")
		d.Usage()
	} else {
		switch os.Args[1] {
		case "describe":
			d.ParseFlags(os.Args[2:]).GetExperiment().PrintAnalysis()
		default:
			fmt.Fprintln(stderr, "expected 'describe' subcommand")
			d.Usage()
		}
	}
}
