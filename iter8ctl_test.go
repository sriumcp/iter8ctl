package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
	"github.com/stretchr/testify/assert"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	fake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// Mocking myOS.Exit function
type myOSMock struct{}

func (m myOSMock) Exit(code int) {
	if code > 0 {
		panic(fmt.Sprintf("Exiting with error code %v", code))
	}
}

// initTestOS registers mock OS lib
func initTestOS() {
	osExiter = myOSMock{}
}

// initTestExperiment registers fake k8s client with access to a pre-built experiment
func initTestExperiment(experiment *v2alpha1.Experiment) {
	// mock k8s client
	crScheme := k8sruntime.NewScheme()
	err := v2alpha1.AddToScheme(crScheme)
	if err != nil {
		panic("Error while adding to v1alpha1 to new scheme")
	}

	getK8sClient = func(d *DescribeCmd) (runtimeclient.Client, error) {
		fakeClient := fake.NewFakeClientWithScheme(crScheme)
		_ = fakeClient.Create(context.Background(), experiment)
		return fakeClient, nil
	}
}

// initTest registers mock OS lib and fake k8s client
func initTest() {
	initTestOS()
	experiment := v2alpha1.NewExperiment("myexp", "myns").
		WithTarget("target").
		WithStrategy(v2alpha1.StrategyTypeCanary).
		WithRequestCount("request-count").
		Build()
	initTestExperiment(experiment)
}

// Unit tests
// Every invocation of iter8ctl in this test should succeed; specifically, main() should not invoke os.Exit(1)
func TestIter8ctl(t *testing.T) {
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args: []string{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns"},
		},
		{
			Args: []string{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "apiVersion", "v2alpha1"},
		},
	} {
		initTest()
		os.Args = test.Args
		assert.NotPanics(t, func() { main() })
	}
}

func TestIter8ctlInvalidSubcommand(t *testing.T) {
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args: []string{"./iter8ctl"},
		},
		{
			Args: []string{"./iter8ctl", "invalid"},
		},
	} {
		initTest()
		os.Args = test.Args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}

func TestIter8ctlParseError(t *testing.T) {
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args: []string{"./iter8ctl", "describe", "--hello", "world"},
		},
	} {
		initTest()
		os.Args = test.Args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}
func TestIter8ctlInvalidNames(t *testing.T) {
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args: []string{"./iter8ctl", "describe", "-name", ""},
		},
		{
			Args: []string{"./iter8ctl", "describe", "-name", "CapitalizedName"},
		},
		{
			Args: []string{"./iter8ctl", "describe", "-name", "namewith*"},
		},
		{
			Args: []string{"./iter8ctl", "describe", "-name", "cindrella", "-namespace", ""},
		},
		{
			Args: []string{"./iter8ctl", "describe", "-name", "cindrella", "-namespace", "Americano"},
		},
	} {
		initTest()
		os.Args = test.Args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}

func TestIter8ctlInvalidAPIVersion(t *testing.T) {
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args: []string{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "--apiVersion", "v2alpha2"},
		},
	} {
		initTest()
		os.Args = test.Args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
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
