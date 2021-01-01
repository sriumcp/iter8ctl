package iter8ctl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
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

// initStdinWithString populates a byte buffer with input string and injects the buffer as stdin
func initStdinWithString(str string) {
	buffer := bytes.Buffer{}
	buffer.Write([]byte(str))
	stdin = &buffer
}

// redirectStdouterr discards stdout and stderr output if LOG_LEVEL is not well-defined
func redirectStdouterr() {
	_, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		stdout = ioutil.Discard
		stderr = ioutil.Discard
	}
}

func TestInvalidSubcommand(t *testing.T) {
	initTestOS()
	redirectStdouterr()
	for _, args := range [][]string{
		{"./iter8ctl"},
		{"./iter8ctl", "invalid"},
	} {
		os.Args = args
		assert.PanicsWithValue(t, "Exiting with non-zero error code", func() { main() })
	}
}

func TestInvalidExperimentYAML(t *testing.T) {
	initTestOS()
	initStdinWithString("abc")
	redirectStdouterr()
	os.Args = []string{"./iter8ctl", "describe", "-f", "-"}
	assert.PanicsWithValue(t, "Exiting with non-zero error code", func() { main() })
}

func TestNormal(t *testing.T) {
	initTestOS()
	// redirectStdouterr()
	for i := 1; i <= 9; i++ {
		_, testFilename, _, _ := runtime.Caller(0)
		expFilename := fmt.Sprintf("experiment%v.yaml", i)
		expFilepath := filepath.Join(filepath.Dir(testFilename), "testdata", expFilename)
		os.Args = []string{"./iter8ctl", "describe", "-f", expFilepath}
		assert.NotPanics(t, func() { main() })
	}
}
