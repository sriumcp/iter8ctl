package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
	log "github.com/sirupsen/logrus"
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
	log.SetReportCaller(true)
}

// DescribeCmd allows for building up a describe command in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type DescribeCmd struct {
	flagSet             *flag.FlagSet
	experimentName      *string
	experimentNamespace *string
	apiVersion          *string
	kubeconfig          *string
	client              client.Client
	experiment          *v2alpha1.Experiment
	logLevel            *string
	err                 error
}

// describeBuilder returns an initial DescribeCmd struct pointer.
func describeBuilder() *DescribeCmd {
	flagSet := flag.NewFlagSet("describe", flag.ExitOnError)
	experimentName := flagSet.String("name", "", "experiment name")
	experimentNamespace := flagSet.String("namespace", "default", "experiment namespace")
	apiVersion := flagSet.String("apiVersion", "v2alpha1", "experiment api version")
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flagSet.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flagSet.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	return &DescribeCmd{
		flagSet:             flagSet,
		experimentName:      experimentName,
		experimentNamespace: experimentNamespace,
		apiVersion:          apiVersion,
		kubeconfig:          kubeconfig,
	}
}

// parseArgs populates experimentName, experimentNamespace, apiVersion, and kubeconfig variables
func (d *DescribeCmd) parseArgs(args []string) *DescribeCmd {
	if d.err != nil {
		return d
	}
	// d.logLevel = flag.String("logLevel", "none", "log level: none/panic/fatal/error/warn/info/debug/trace")
	d.err = d.flagSet.Parse(args)
	if d.err != nil {
		log.WithFields(log.Fields{
			"msg": d.err,
		}).Error()
	}
	return d
}

// validateName validates experiment name
func (d *DescribeCmd) validateName() *DescribeCmd {
	if d.err != nil {
		return d
	}

	var namePrefix *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])`)
	var nameSuffix *regexp.Regexp = regexp.MustCompile(`([[:lower:]]|[[:digit:]])$`)
	var nameRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]]|-|\.){1,253}`)
	if !(namePrefix.MatchString(*d.experimentName) && nameSuffix.MatchString(*d.experimentName) && nameRegex.MatchString(*d.experimentName)) {
		errMsg := "Invalid experiment name... name should contain no more than 253 characters and only lowercase alphanumeric characters, '-' or '.'... name should start and end with an alphanumeric character"
		d.err = errors.New(errMsg)
		fmt.Println(errMsg)
	}

	if d.err != nil {
		log.WithFields(log.Fields{
			"msg": d.err,
		}).Error()
	}
	return d
}

// validateNamespace validates experiment namespace
func (d *DescribeCmd) validateNamespace() *DescribeCmd {
	if d.err != nil {
		return d
	}

	var namespacePrefix *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])`)
	var namespaceSuffix *regexp.Regexp = regexp.MustCompile(`([[:lower:]]|[[:digit:]])$`)
	var namespaceRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]]|-){1,63}`)
	if !(namespacePrefix.MatchString(*d.experimentNamespace) && namespaceSuffix.MatchString(*d.experimentNamespace) && namespaceRegex.MatchString(*d.experimentNamespace)) {
		errMsg := "Invalid experiment namespace... namespace should contain no more than 63 characters and only lowercase alphanumeric characters, or '-'... namespace should start and end with an alphanumeric character"
		d.err = errors.New(errMsg)
		fmt.Println(errMsg)
	}

	if d.err != nil {
		log.WithFields(log.Fields{
			"msg": d.err,
		}).Error()
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
		errMsg := "Invalid experiment APIVersion... only allowed value for APIVersion is 'v2alpha1'"
		d.err = errors.New(errMsg)
		fmt.Println(errMsg)
	}

	if d.err != nil {
		log.WithFields(log.Fields{
			"msg": d.err,
		}).Error()
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

var getK8sClient = func(d *DescribeCmd) (runtimeclient.Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *d.kubeconfig)
	if err != nil {
		return nil, err
	}

	crScheme := runtime.NewScheme()
	err = v2alpha1.AddToScheme(crScheme)
	if err != nil {
		return nil, err
	}

	rc, err := runtimeclient.New(config, client.Options{
		Scheme: crScheme,
	})
	if err != nil {
		return nil, err
	}
	return rc, nil
}

// setK8sClient sets the clientset variable within DescribeCmd struct
func (d *DescribeCmd) setK8sClient() *DescribeCmd {
	if d.err != nil {
		return d
	}

	d.client, d.err = getK8sClient(d)
	if d.err != nil {
		log.WithFields(log.Fields{
			"msg": d.err,
		}).Error()
	}
	return d
}

// getExperiment gets the experiment resource object from the k8s cluster
func (d *DescribeCmd) getExperiment() *DescribeCmd {
	if d.err != nil {
		return d
	}
	d.experiment = &v2alpha1.Experiment{}
	d.err = d.client.Get(context.Background(), client.ObjectKey{
		Namespace: *d.experimentNamespace,
		Name:      *d.experimentName,
	}, d.experiment)
	if d.err != nil {
		log.WithFields(log.Fields{
			"msg": d.err,
		}).Error()
		fmt.Println("Cannot get experiment object. Ensure you supplied valid arguments to 'iter8ctl describe' command")
	} else {
		data, _ := json.MarshalIndent(d.experiment, "", "  ")
		log.Info("\nGot experiment...\n", string(data))
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
		d.parseArgs(os.Args[2:]).validate().setK8sClient().getExperiment().printAnalysis()
		if d.err != nil {
			d.flagSet.Usage()
		}

	default:
		fmt.Println("expected 'describe' subcommand")
		osExiter.Exit(1)
	}
}
