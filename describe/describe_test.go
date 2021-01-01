package describe

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/iter8-tools/iter8ctl/utils"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
	Cmd := Builder(&a, &b, &c)
	Cmd.Usage()
	assert.Greater(t, c.Len(), 0)
	assert.NoError(t, Cmd.Error())
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		{"-name", "helloworld"},
		{},
		{"-f"},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs(args)
		assert.Error(t, Cmd.Error())
		d := Cmd.Error()
		Cmd.GetExperiment()
		assert.Error(t, Cmd.Error())
		assert.Equal(t, d, Cmd.Error())
	}
}

func TestFileInputGood(t *testing.T) {
	for _, args := range [][]string{
		{"-f", utils.CompletePath("../", "testdata/experiment1.yaml")},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs(args).GetExperiment()
		assert.NoError(t, Cmd.Error())
	}
}

func TestStdinGood(t *testing.T) {
	for _, args := range [][]string{
		{"-f", "-"},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		data, _ := ioutil.ReadFile(utils.CompletePath("../", "testdata/experiment1.yaml"))
		a.Write(data)
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs(args).GetExperiment()
		assert.NoError(t, Cmd.Error())
	}
}

func TestFileDoesNotExist(t *testing.T) {
	for _, args := range [][]string{
		{"-f", utils.CompletePath("../", "testdata/nonexistant.yaml")},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs(args).GetExperiment()
		assert.Error(t, Cmd.Error())
	}
}

func TestStdinBadYAML(t *testing.T) {
	for _, args := range [][]string{
		{"-f", "-", "abc 123 xyz"},
		{"-f", "-", "abc:\nabc:1"},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		a.Write([]byte(args[2]))
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs(args[:2]).GetExperiment()
		assert.Error(t, Cmd.Error())
	}
}

func TestPrintProgress(t *testing.T) {
	for i := 1; i <= 9; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().PrintProgress()
		assert.NoError(t, Cmd.Error())
	}
}

func TestPrintWinnerAssessment(t *testing.T) {
	for i := 1; i <= 9; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().PrintWinnerAssessment()
		assert.NoError(t, Cmd.Error())
	}
}

func TestPrintObjectiveAssessment(t *testing.T) {
	for i := 1; i <= 9; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().PrintObjectiveAssessment()
		assert.NoError(t, Cmd.Error())
	}
}

func TestPrintVersionAssessment(t *testing.T) {
	for i := 1; i <= 9; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().PrintVersionAssessment()
		assert.NoError(t, Cmd.Error())
	}
}

func TestPrintMetrics(t *testing.T) {
	for i := 1; i <= 9; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().PrintMetrics()
		assert.NoError(t, Cmd.Error())
	}
}

func TestPrintAnalysis(t *testing.T) {
	for i := 1; i <= 9; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().PrintAnalysis()
		assert.NoError(t, Cmd.Error())
	}
}
