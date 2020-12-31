package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		panic("Exiting with non-zero error code")
	}
}

// initTestOS registers the mock OS struct (testOS) defined above
func initTestOS() {
	osExiter = &testOS{}
}

// initStdinWithFile populates a byte buffer with input file contents and injects the buffer as stdin
func initStdinWithFile(filePath string) {
	data, _ := ioutil.ReadFile(filePath)
	buffer := bytes.Buffer{}
	buffer.Write(data)
	stdin = &buffer
}

// initStdinWithFile populates a byte buffer with input string and injects the buffer as stdin
func initStdinWithString(str string) {
	buffer := bytes.Buffer{}
	buffer.Write([]byte(str))
	stdin = &buffer
}

/* Unit tests */

func TestFilepath(t *testing.T) {
	initTestOS()
	for i := 1; i <= 8; i++ {
		_, testFilename, _, _ := runtime.Caller(0)
		expFilename := fmt.Sprintf("experiment%v.yaml", i)
		expFilepath := filepath.Join(filepath.Dir(testFilename), "testdata", expFilename)
		os.Args = []string{"./iter8ctl", "describe", "-f", expFilepath}
		assert.NotPanics(t, func() { main() })
	}
}

func TestStdin(t *testing.T) {
	initTestOS()
	for i := 1; i <= 8; i++ {
		_, testFilename, _, _ := runtime.Caller(0)
		expFilename := fmt.Sprintf("experiment%v.yaml", i)
		expFilepath := filepath.Join(filepath.Dir(testFilename), "testdata", expFilename)
		initStdinWithFile(expFilepath)
		os.Args = []string{"./iter8ctl", "describe", "-f", "-"}
		assert.NotPanics(t, func() { main() })
	}
}

func TestInvalidSubcommand(t *testing.T) {
	initTestOS()
	for _, args := range [][]string{
		{"./iter8ctl"},
		{"./iter8ctl", "invalid"},
	} {
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with non-zero error code", func() { main() })
	}
}

func TestParseError(t *testing.T) {
	initTestOS()
	os.Args = []string{"./iter8ctl", "describe", "--hello", "world"}
	assert.PanicsWithValue(t, "Exiting with non-zero error code", func() { main() })
}

func TestInvalidYAML(t *testing.T) {
	initTestOS()
	initStdinWithString("playing_playlist: {{ action }} playlist {{ playlist_name }}")
	os.Args = []string{"./iter8ctl", "describe", "-f", "-"}
	assert.PanicsWithValue(t, "Exiting with non-zero error code", func() { main() })
}

func TestInvalidFile(t *testing.T) {
	initTestOS()
	os.Args = []string{"./iter8ctl", "describe", "-f", "abc123xyz789.yaml.json"}
	assert.PanicsWithValue(t, "Exiting with non-zero error code", func() { main() })
}

func TestInvalidExperimentYAML(t *testing.T) {
	initTestOS()
	initStdinWithString("abc")
	os.Args = []string{"./iter8ctl", "describe", "-f", "-"}
	assert.PanicsWithValue(t, "Exiting with non-zero error code", func() { main() })
}

func TestPrintAnalysis(t *testing.T) {
	initTestOS()
	for i := 1; i <= 9; i++ {
		_, testFilename, _, _ := runtime.Caller(0)
		expFilename := fmt.Sprintf("experiment%v.yaml", i)
		expFilepath := filepath.Join(filepath.Dir(testFilename), "testdata", expFilename)
		os.Args = []string{"./iter8ctl", "describe", "-f", expFilepath}
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
