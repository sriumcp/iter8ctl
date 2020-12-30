package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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

/* Unit tests */

func TestDescribeBuilder(t *testing.T) {
	initTestOS()
	for i := 1; i <= 8; i++ {
		_, testFilename, _, _ := runtime.Caller(0)
		expFilename := fmt.Sprintf("experiment%v.yaml", i)
		expFilepath := filepath.Join(filepath.Dir(testFilename), "testdata", expFilename)
		os.Args = []string{"./iter8ctl", "describe", "-e", expFilepath}
		assert.NotPanics(t, func() { main() })
	}
}

func TestInvalidSubcommand(t *testing.T) {
	initTestOS()
	for _, args := range [][]string{
		{"./iter8ctl"}, {"./iter8ctl", "invalid"},
	} {
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}

func TestParseError(t *testing.T) {
	initTestOS()
	for _, args := range [][]string{
		{"./iter8ctl", "describe", "--hello", "world"},
	} {
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with error code 1", func() { main() })
	}
}

func TestPrintAnalysis(t *testing.T) {
	initTestOS()
	for i := 1; i <= 8; i++ {
		_, testFilename, _, _ := runtime.Caller(0)
		expFilename := fmt.Sprintf("experiment%v.yaml", i)
		expFilepath := filepath.Join(filepath.Dir(testFilename), "testdata", expFilename)
		os.Args = []string{"./iter8ctl", "describe", "-e", expFilepath}
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
