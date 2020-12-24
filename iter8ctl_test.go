package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestIter8ctl(t *testing.T) {
	for _, test := range []struct {
		Args   []string
		Output string
	}{
		{
			Args:   []string{"./iter8ctl", "describe", "--name", "myexp"},
			Output: "",
		},
		{
			Args:   []string{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns"},
			Output: "",
		},
		{
			Args:   []string{"./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "apiVersion", "v2alpha1"},
			Output: "",
		},
	} {
		t.Run("", func(t *testing.T) {
			os.Args = test.Args
			out = bytes.NewBuffer(nil)
			main()

			if actual := out.(*bytes.Buffer).String(); actual != test.Output {
				t.Errorf("expected %s, but got %s", test.Output, actual)
			}
		})
	}
}

func execExitingCommand(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	actualOutput := fmt.Sprintf("%s", &out)
	return actualOutput, err
}

func postProcess(t *testing.T, actualOutput string, expectedOutput string, err error) {
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if strings.HasPrefix(actualOutput, expectedOutput) {
			return
		}
		t.Logf("Expected prefix for output: %s", expectedOutput)
		t.Fatalf("Actual output: %s", actualOutput)
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestIter8ctlInvalidSubcommand(t *testing.T) {
	expectedOutput := "expected 'describe' subcommand"

	cmd := exec.Command("./iter8ctl")
	actualOutput, err := execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)

	cmd = exec.Command("./iter8ctl", "invalid")
	actualOutput, err = execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)
}

func TestIter8ctlInvalidNames(t *testing.T) {
	expectedOutput := "Expected a valid value for (experiment) name and namespace."

	cmd := exec.Command("./iter8ctl", "describe", "-name", "")
	actualOutput, err := execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)

	cmd = exec.Command("./iter8ctl", "describe", "-name", "CapitalizedName")
	actualOutput, err = execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)

	cmd = exec.Command("./iter8ctl", "describe", "-name", "namewith*")
	actualOutput, err = execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)

	cmd = exec.Command("./iter8ctl", "describe", "-name", "cindrella", "-namespace", "")
	actualOutput, err = execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)

	cmd = exec.Command("./iter8ctl", "describe", "-name", "cindrella", "-namespace", "Americano")
	actualOutput, err = execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)
}

func TestIter8ctlInvalidAPIVersion(t *testing.T) {
	expectedOutput := "Expected a valid value for (experiment) api version."

	cmd := exec.Command("./iter8ctl", "describe", "--name", "myexp", "--namespace", "myns", "--apiVersion", "v2alpha2")
	actualOutput, err := execExitingCommand(cmd)
	postProcess(t, actualOutput, expectedOutput, err)
}
