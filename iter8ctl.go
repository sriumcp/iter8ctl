package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
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

// K8sClient interface enables getting a k8s client
type K8sClient interface {
	getK8sClient(kubeconfigPath *string) (runtimeclient.Client, error)
}
type iter8ctlK8sClient struct{}

func (k iter8ctlK8sClient) getK8sClient(kubeconfigPath *string) (runtimeclient.Client, error) {
	crScheme := runtime.NewScheme()
	err := v2alpha1.AddToScheme(crScheme)
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
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

var k8sClient K8sClient

// init initializes logging, osExiter, and k8sClient
func init() {
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
	osExiter = &iter8ctlOS{}
	k8sClient = &iter8ctlK8sClient{}
}

// DescribeCmd allows for building up a describe command in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type DescribeCmd struct {
	K8sClient
	flagSet             *flag.FlagSet
	experimentName      *string
	experimentNamespace *string
	apiVersion          *string
	kubeconfigPath      *string
	client              client.Client
	experiment          *v2alpha1.Experiment
	err                 error
}

// describeBuilder returns an initial DescribeCmd struct pointer.
func describeBuilder(k K8sClient) *DescribeCmd {
	flagSet := flag.NewFlagSet("describe", flag.ContinueOnError)
	experimentName := flagSet.String("name", "", "experiment name")
	experimentNamespace := flagSet.String("namespace", "default", "experiment namespace")
	apiVersion := flagSet.String("apiVersion", "v2alpha1", "experiment api version")
	var kubeconfigPath *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = flagSet.String("kubeconfigPath", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfigPath = flagSet.String("kubeconfigPath", "", "absolute path to the kubeconfig file")
	}

	return &DescribeCmd{
		K8sClient:           k,
		flagSet:             flagSet,
		experimentName:      experimentName,
		experimentNamespace: experimentNamespace,
		apiVersion:          apiVersion,
		kubeconfigPath:      kubeconfigPath,
		client:              nil,
		experiment:          &v2alpha1.Experiment{},
		err:                 nil,
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

// validateName validates experimentName
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

	return d
}

// validateNamespace validates experimentNamespace
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

	return d
}

// validate validates experimentName, experimentNamespace, and apiVersion
func (d *DescribeCmd) validate() *DescribeCmd {
	return d.validateName().validateNamespace().validateAPIVersion()
}

// setK8sClient sets the clientset variable within DescribeCmd struct
func (d *DescribeCmd) setK8sClient() *DescribeCmd {
	if d.err != nil {
		return d
	}
	d.client, d.err = d.getK8sClient(d.kubeconfigPath)
	if d.err != nil {
		fmt.Printf("Error setting k8s client: %s\n", d.err)
	}
	return d
}

// getExperiment gets the experiment resource object from the k8s cluster
func (d *DescribeCmd) getExperiment() *DescribeCmd {
	if d.err != nil {
		return d
	}
	ro := v2alpha1.ExperimentList{
		TypeMeta: v1.TypeMeta{},
		ListMeta: v1.ListMeta{},
		Items:    []v2alpha1.Experiment{},
	}
	d.client.List(context.Background(), &ro)
	d.experiment = &v2alpha1.Experiment{}
	d.err = d.client.Get(context.Background(), client.ObjectKey{
		Namespace: *d.experimentNamespace,
		Name:      *d.experimentName,
	}, d.experiment)
	if d.err != nil {
		fmt.Printf("Cannot get experiment object. %s\n", d.err)
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
		d := describeBuilder(k8sClient)
		d.parseArgs(os.Args[2:]).validate().setK8sClient().getExperiment().printAnalysis()
		if d.err != nil {
			osExiter.Exit(1)
		}

	default:
		fmt.Println("expected 'describe' subcommand")
		osExiter.Exit(1)
	}
}
