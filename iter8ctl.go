package main

import (
	"flag"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/iter8-tools/iter8ctl/exp"
)

// Rules for valid k8s resource name and namespace are here: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
var nameRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])([[:lower:]]|[[:digit:]]|-|\.){0,253}([[:lower:]]|[[:digit:]])$`)
var namespaceRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])([[:lower:]]|[[:digit:]]|-){0,63}([[:lower:]]|[[:digit:]])$`)

// Exact match required for apiVersion: v2alpha1 is the only allowed value
var apiVersionRegex *regexp.Regexp = regexp.MustCompile(`\bv2alpha1\b`)

// OSExiter wraps os.Exit(1) calls. Useful for mocks in unit tests.
// Reference: https://medium.com/@ankur_anand/how-to-mock-in-your-go-golang-tests-b9eee7d7c266
type OSExiter interface {
	Exit(code int)
}

type myOS struct{}

func (m myOS) Exit(code int) {
	os.Exit(code)
}

var osExiter OSExiter

func init() {
	// Initializing exiter
	osExiter = myOS{}

	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// With the default log.SetFormatter(&log.TextFormatter{}) when a TTY is not attached, the output is compatible with the logfmt format
	// Above comment from here: https://github.com/sirupsen/logrus
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {

	describeCmd := flag.NewFlagSet("describe", flag.ExitOnError)
	experimentName := describeCmd.String("name", "", "experiment name")
	experimentNamespace := describeCmd.String("namespace", "default", "experiment namespace")
	apiVersion := describeCmd.String("apiVersion", "v2alpha1", "experiment api version")

	if len(os.Args) < 2 {
		log.Error("expected 'describe' subcommand")
		osExiter.Exit(1)
	}

	switch os.Args[1] {

	case "describe":
		describeCmd.Parse(os.Args[2:])

		if len(*experimentName) == 0 || len(*experimentName) > 253 || !nameRegex.MatchString(*experimentName) || len(*experimentNamespace) == 0 || len(*experimentNamespace) > 63 || !namespaceRegex.MatchString(*experimentNamespace) {
			log.WithFields(log.Fields{
				"name":      *experimentName,
				"namespace": *experimentNamespace,
			}).Error("expected a valid value for (experiment) name and namespace.")
			log.Error("name should contain no more than 253 characters and only lowercase alphanumeric characters, '-' or '.'; start and end with an alphanumeric character.")
			log.Error("namespace should contain no more than 63 characters and only lowercase alphanumeric characters, or '-'; start and end with an alphanumeric character.")
			osExiter.Exit(1)
		}

		if !apiVersionRegex.MatchString(*apiVersion) {
			log.WithFields(log.Fields{
				"apiVersion": *apiVersion,
			}).Error("expected a valid value for (experiment) api version.")
			osExiter.Exit(1)
		}

		experiment, err := exp.GetExperiment(experimentName, experimentNamespace, apiVersion)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Encountered error while getting experiment")
			osExiter.Exit(1)
		}

		experiment.PrintAnalysis()

	default:
		log.WithFields(log.Fields{
			"subcommand": os.Args[1],
		}).Error("expected 'describe' subcommand")
		osExiter.Exit(1)
	}
}
