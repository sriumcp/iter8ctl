package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"regexp"

	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/node/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// OSExiter wraps os.Exit(1) calls. Useful for mocks in unit tests.
type OSExiter interface {
	Exit(code int)
}
type myOS struct{}

func (m myOS) Exit(code int) {
	os.Exit(code)
}

var osExiter OSExiter

// init initializes osExiter and logging
func init() {
	osExiter = myOS{}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(log.WarnLevel)
}

// DescribeCmd allows for building up a describe command in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type DescribeCmd struct {
	experimentName      *string
	experimentNamespace *string
	apiVersion          *string
	kubeconfig          *string
	client              client.Client
	experiment          *v2alpha1.Experiment
	err                 error
}

// describeBuilder returns a DescribeCmd struct pointer with struct variables initialized to 'nil' values.
func describeBuilder() *DescribeCmd {
	return &DescribeCmd{}
}

// parseArgs populates experimentName, experimentNamespace, apiVersion, and kubeconfig variables
func (d *DescribeCmd) parseArgs(args []string) *DescribeCmd {
	if d.err != nil {
		return d
	}

	describeCmd := flag.NewFlagSet("describe", flag.ContinueOnError)
	d.experimentName = describeCmd.String("name", "", "experiment name")
	d.experimentNamespace = describeCmd.String("namespace", "default", "experiment namespace")
	d.apiVersion = describeCmd.String("apiVersion", "v2alpha1", "experiment api version")
	if home := homedir.HomeDir(); home != "" {
		d.kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		d.kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	d.err = describeCmd.Parse(args)

	return d
}

// validateName validates experiment name
func (d *DescribeCmd) validateName() *DescribeCmd {
	if d.err != nil {
		return d
	}

	var nameRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])([[:lower:]]|[[:digit:]]|-|\.){0,253}([[:lower:]]|[[:digit:]])$`)
	if len(*d.experimentName) == 0 || len(*d.experimentName) > 253 || !nameRegex.MatchString(*d.experimentName) {
		d.err = errors.New("Invalid experiment name; name should contain no more than 253 characters and only lowercase alphanumeric characters, '-' or '.'; start and end with an alphanumeric character")
	}

	return d
}

// validateNamespace validates experiment namespace
func (d *DescribeCmd) validateNamespace() *DescribeCmd {
	if d.err != nil {
		return d
	}

	var namespaceRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])([[:lower:]]|[[:digit:]]|-){0,63}([[:lower:]]|[[:digit:]])$`)
	if len(*d.experimentNamespace) == 0 || len(*d.experimentNamespace) > 63 || !namespaceRegex.MatchString(*d.experimentNamespace) {
		d.err = errors.New("Invalid experiment namespace; namespace should contain no more than 63 characters and only lowercase alphanumeric characters, or '-'; start and end with an alphanumeric character")
	}

	return d
}

// validateAPIVersion validates apiVersion
func (d *DescribeCmd) validateAPIVersion() *DescribeCmd {
	if d.err != nil {
		return d
	}

	var apiVersionRegex *regexp.Regexp = regexp.MustCompile(`\bv2alpha1\b`)
	if !apiVersionRegex.MatchString(*d.apiVersion) {
		d.err = errors.New("Invalid experiment APIVersion; only allowed value for APIVersion is 'v2alpha1'")
	}

	return d
}

// validate validates experimentName, experimentNamespace, and apiVersion
func (d *DescribeCmd) validate() *DescribeCmd {
	if d.err != nil {
		return d
	}

	return d.validateName().validateNamespace().validateAPIVersion()
}

// setK8sClient sets the clientset variable within DescribeCmd struct
func (d *DescribeCmd) setK8sClient() *DescribeCmd {
	if d.err != nil {
		return d
	}

	config, err := clientcmd.BuildConfigFromFlags("", *d.kubeconfig)
	if err != nil {
		d.err = err
		return d
	}
	crScheme := runtime.NewScheme()

	err = v1alpha1.AddToScheme(crScheme)
	if err != nil {
		d.err = err
		return d
	}

	d.client, err = runtimeclient.New(config, client.Options{
		Scheme: crScheme,
	})
	if err != nil {
		d.err = err
		return d
	}

	return d
}

// getExperiment gets the experiment resource object from the k8s cluster
func (d *DescribeCmd) getExperiment() *DescribeCmd {
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
		log.Error("expected 'describe' subcommand")
		osExiter.Exit(1)
	}

	switch os.Args[1] {

	case "describe":
		desc := describeBuilder()
		err := desc.parseArgs(os.Args[2:]).validate().setK8sClient().getExperiment().printAnalysis()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("'describe' command resulted in error")
			osExiter.Exit(1)
		}

	default:
		log.WithFields(log.Fields{
			"subcommand": os.Args[1],
		}).Error("expected 'describe' subcommand")
		osExiter.Exit(1)
	}
}
