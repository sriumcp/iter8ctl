package iter8ctl

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	"github.com/iter8-tools/iter8ctl/utils"
	"github.com/stretchr/testify/assert"
)

/* CLI tests */

func TestCLI(t *testing.T) {
	// All subtests within this test rely in ./iter8ctl. So go build.
	exec.Command("go", "build", utils.CompletePath("", "")).Run()

	type test struct {
		name           string   // name of this test
		flags          []string // flags supplied to .iter8ctl command
		outputFilename string   // relative to testdata
		errorFilename  string   // relative to testdata
	}

	tests := []test{
		{name: "no-flags", flags: []string{}, outputFilename: "", errorFilename: "error-no-flags.txt"},
		{name: "invalid-subcommand", flags: []string{"invalid"}, outputFilename: "", errorFilename: "error-invalid-subcommand.txt"},
		{name: "undefined-flag", flags: []string{"describe", "-name", "helloworld"}, outputFilename: "", errorFilename: "error-undefined-flag.txt"},
		{name: "experiment1", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment1.yaml")}, outputFilename: "experiment1.out", errorFilename: ""},
		{name: "experiment2", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment2.yaml")}, outputFilename: "experiment2.out", errorFilename: ""},
		{name: "experiment3", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment3.yaml")}, outputFilename: "experiment3.out", errorFilename: ""},
		{name: "experiment4", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment4.yaml")}, outputFilename: "experiment4.out", errorFilename: ""},
		{name: "experiment5", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment5.yaml")}, outputFilename: "experiment5.out", errorFilename: ""},
		{name: "experiment6", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment6.yaml")}, outputFilename: "experiment6.out", errorFilename: ""},
		{name: "experiment7", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment7.yaml")}, outputFilename: "experiment7.out", errorFilename: ""},
		{name: "experiment8", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment8.yaml")}, outputFilename: "experiment8.out", errorFilename: ""},
		{name: "experiment9", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment9.yaml")}, outputFilename: "experiment9.out", errorFilename: ""},
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

			if tc.errorFilename != "" {
				assert.Error(t, cmd.Run())
				t.Log("Reading error file: ", utils.CompletePath("testdata", tc.errorFilename))
				b3, err := ioutil.ReadFile(utils.CompletePath("testdata", tc.errorFilename))
				if err != nil {
					t.Fatal("Unable to read error file contents")
				}
				assert.Equal(t, string(b3), b1.String())
			} else {
				assert.NoError(t, cmd.Run())
			}

			if tc.outputFilename != "" {
				t.Log("Reading output file: ", utils.CompletePath("testdata", tc.outputFilename))
				b4, err := ioutil.ReadFile(utils.CompletePath("testdata", tc.outputFilename))
				if err != nil {
					t.Fatal("Unable to read output file contents")
				}
				assert.Equal(t, string(b4), b2.String())
			}
		})
	}
}
