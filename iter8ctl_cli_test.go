package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// pointer is a helper function for creating string pointers from literals
func pointer(s string) *string {
	return &s
}

// readFileAsBytes is a helper function for reading contents of a file as []byte; filename is relative to testdata folder
func readFileAsBytes(filename string) ([]byte, error) {
	_, testFilename, _, _ := runtime.Caller(0)
	filePath := filepath.Join(filepath.Dir(testFilename), "testdata", filename)
	return ioutil.ReadFile(filePath)
}

/* CLI tests */

func TestCLI(t *testing.T) {
	// All subtests within this test rely in ./iter8ctl. So go build if needed.
	if _, err := os.Stat("iter8ctl"); os.IsNotExist(err) {
		exec.Command("go", "build").Run()
	}

	type test struct {
		name           string   // name of this test
		flags          []string // flags supplied to .iter8ctl command
		outputFilename *string  // relative to testdata
		errorFilename  *string  // relative to testdata
	}

	tests := []test{
		{name: "no-flags", flags: []string{"./iter8ctl"}, outputFilename: nil, errorFilename: pointer("error-no-flags.txt")},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cmd := exec.Command("./iter8ctl", tc.flags...)
			b1 := &strings.Builder{}
			cmd.Stderr = b1
			b2 := &strings.Builder{}
			cmd.Stdout = b2

			if tc.errorFilename != nil {
				assert.Error(t, cmd.Run())
				b3, err := readFileAsBytes(*tc.errorFilename)
				if err != nil {
					t.Fatal("Unable to read error file contents")
				}
				assert.Equal(t, string(b3), b1.String())
			} else {
				assert.NoError(t, cmd.Run())
			}

			if tc.outputFilename != nil {
				b4, err := readFileAsBytes(*tc.outputFilename)
				if err != nil {
					t.Fatal("Unable to read output file contents")
				}
				assert.Equal(t, string(b4), b1.String())
			}
		})
	}
}

// func Test_ExecuteCommand(t *testing.T) {
// 	cmd := NewRootCmd("hi")
// 	b := bytes.NewBufferString("")
// 	cmd.SetOut(b)
// 	cmd.Execute()
// 	out, err := ioutil.ReadAll(b)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if string(out) != "hi" {
// 		t.Fatalf("expected \"%s\" got \"%s\"", "hi", string(out))
// 	}
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
