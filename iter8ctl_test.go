package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
	"github.com/stretchr/testify/assert"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	fake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"
)

// Mocking os.Exit function
type testOS struct{}

func (t *testOS) Exit(code int) {
	if code > 0 {
		panic(fmt.Sprintf("Exiting with error code %v", code))
	}
}

// initTestOS registers the mock OS struct (testOS) defined above
func initTestOS() {
	osExiter = &testOS{}
}

// Mocking iter8ctlK8sClient.getK8sClient function
type testK8sClient struct {
	experiment *v2alpha1.Experiment
}

func (t *testK8sClient) getK8sClient(kubeconfigPath *string) (runtimeclient.Client, error) {
	crScheme := k8sruntime.NewScheme()
	err := v2alpha1.AddToScheme(crScheme)
	if err != nil {
		panic("Error while adding to v1alpha1 to new scheme")
	}
	if t.experiment == nil {
		return fake.NewFakeClientWithScheme(crScheme), nil
	}
	return fake.NewFakeClientWithScheme(crScheme, t.experiment), nil
}

// initTestK8sClient enables a fake k8s client with access to a pre-built experiment
func initTestK8sClient(experiment *v2alpha1.Experiment) {
	k8sClient = &testK8sClient{experiment}
}

// initTest registers mock OS lib and fake k8s client
func initTestWithMyExp() {
	initTestOS()
	experiment := v2alpha1.NewExperiment("myexp", "myns").
		WithTarget("target").
		WithStrategy(v2alpha1.StrategyTypeCanary).
		WithRequestCount("request-count").
		Build()
	initTestK8sClient(experiment)
}

/* Unit tests */

func TestDescribeBuilder(t *testing.T) {
	for _, args := range [][]string{
		{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns"},
		{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "apiVersion", "v2alpha1"},
	} {
		initTestWithMyExp()
		os.Args = args
		assert.NotPanics(t, func() { main() })
	}
}

func TestInvalidSubcommand(t *testing.T) {
	for _, args := range [][]string{
		{"./iter8ctl"}, {"./iter8ctl", "invalid"},
	} {
		initTestWithMyExp()
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}

func TestParseError(t *testing.T) {
	for _, args := range [][]string{
		{"./iter8ctl", "describe", "--hello", "world"},
	} {
		initTestWithMyExp()
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}
func TestInvalidNames(t *testing.T) {
	for _, args := range [][]string{
		{"./iter8ctl", "describe", "-name", ""},
		{"./iter8ctl", "describe", "-name", "CapitalizedName"},
		{"./iter8ctl", "describe", "-name", "namewith*"},
		{"./iter8ctl", "describe", "-name", "cindrella", "-namespace", ""},
		{"./iter8ctl", "describe", "-name", "cindrella", "-namespace", "Americano"},
	} {
		initTestWithMyExp()
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}

func TestInvalidAPIVersion(t *testing.T) {
	for _, args := range [][]string{
		{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "--apiVersion", "v2alpha2"},
	} {
		initTestWithMyExp()
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}

func TestInvalidKubeconfigPath(t *testing.T) {
	_, testFilename, _, _ := runtime.Caller(0)
	kubeconfigPath := filepath.Join(filepath.Dir(testFilename), "testdata", "kubeconfig")
	for _, args := range [][]string{
		{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "--apiVersion", "v2alpha1", "--kubeconfigPath", kubeconfigPath},
	} {
		initTestOS()
		os.Args = args
		d := describeBuilder(&iter8ctlK8sClient{})
		d.parseArgs(os.Args[2:]).validate().setK8sClient()
		assert.Error(t, d.err)
	}
}

func TestPrintAnalysis(t *testing.T) {
	initTestOS()
	for i := 1; i <= 8; i++ {
		_, testFilename, _, _ := runtime.Caller(0)
		expFilename := fmt.Sprintf("experiment%v.yaml", i)
		expFilepath := filepath.Join(filepath.Dir(testFilename), "testdata", expFilename)
		expBytes, _ := ioutil.ReadFile(expFilepath)
		experiment := &v2alpha1.Experiment{}
		err := yaml.Unmarshal(expBytes, experiment)
		if err != nil {
			t.Error(err)
		}
		os.Args = []string{"./iter8ctl", "describe", "--name", "sklearn-iris-experiment-1", "--namespace", "kfserving-test"}
		initTestK8sClient(experiment)
		assert.NotPanics(t, func() { main() })
	}
}

// // Integration tests
// func execExitingCommand(cmd *exec.Cmd) (string, error) {
// 	var out bytes.Buffer
// 	cmd.Stderr = &out
// 	err := cmd.Run()
// 	actualOutput := fmt.Sprintf("%s", &out)
// 	return actualOutput, err
// }

// func postProcess(t *testing.T, actualOutput string, expectedOutput string, err error) {
// 	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
// 		if strings.Contains(actualOutput, expectedOutput) {
// 			return
// 		}
// 		t.Logf("expected substring in output: %s", expectedOutput)
// 		t.Fatalf("actual output: %s", actualOutput)
// 	}
// 	t.Fatalf("process ran with err %v, want exit status 1", err)
// }

// func TestIter8ctlInvalidSubcommandIntegration(t *testing.T) {
// 	expectedOutput := "expected 'describe' subcommand"

// 	cmd := exec.Command("./iter8ctl")
// 	actualOutput, err := execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)

// 	cmd = exec.Command("./iter8ctl", "invalid")
// 	actualOutput, err = execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)
// }

// func TestIter8ctlInvalidNamesIntegration(t *testing.T) {
// 	expectedOutput := "expected a valid value for (experiment) name and namespace."

// 	cmd := exec.Command("./iter8ctl", "describe", "-name", "")
// 	actualOutput, err := execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)

// 	cmd = exec.Command("./iter8ctl", "describe", "-name", "CapitalizedName")
// 	actualOutput, err = execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)

// 	cmd = exec.Command("./iter8ctl", "describe", "-name", "namewith*")
// 	actualOutput, err = execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)

// 	cmd = exec.Command("./iter8ctl", "describe", "-name", "cindrella", "-namespace", "")
// 	actualOutput, err = execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)

// 	cmd = exec.Command("./iter8ctl", "describe", "-name", "cindrella", "-namespace", "Americano")
// 	actualOutput, err = execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)
// }

// func TestIter8ctlInvalidAPIVersionIntegration(t *testing.T) {
// 	expectedOutput := "expected a valid value for (experiment) api version."

// 	cmd := exec.Command("./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "--apiVersion", "v2alpha2")
// 	actualOutput, err := execExitingCommand(cmd)
// 	postProcess(t, actualOutput, expectedOutput, err)
// }
