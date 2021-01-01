package experiment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/iter8-tools/iter8ctl/utils"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

// getExp is a helper function for extracting an experiment object from experiment filenamePrefix
// filePath is relative to testdata folder
func getExp(filenamePrefix string) (*Experiment, error) {
	experimentFilepath := utils.CompletePath("../", fmt.Sprintf("testdata/%s.yaml", filenamePrefix))
	expBytes, err := ioutil.ReadFile(experimentFilepath)
	if err != nil {
		return nil, err
	}
	expBytesJSON, err := yaml.YAMLToJSON(expBytes)
	if err != nil {
		return nil, err
	}
	exp := &Experiment{}
	err = json.Unmarshal(expBytesJSON, &exp)
	if err != nil {
		return nil, err
	}
	return exp, nil
}

type test struct {
	name                   string // name of this test
	started                bool
	exp                    *Experiment
	errorRates, fakeMetric []string
	satisfyStrs, fakeObj   []string
}

var fakeValStrs = []string{"unavailable", "unavailable"}

var satisfyStrs = []string{"true", "true"}

var tests = []test{
	{name: "experiment1", started: false, errorRates: []string{}, fakeMetric: []string{}, satisfyStrs: []string{}, fakeObj: []string{}},
	{name: "experiment2", started: false, errorRates: []string{}, fakeMetric: []string{}, satisfyStrs: []string{}, fakeObj: []string{}},
	{name: "experiment3", started: true, errorRates: []string{"0", "0"}, fakeMetric: fakeValStrs, satisfyStrs: satisfyStrs, fakeObj: fakeValStrs},
	{name: "experiment4", started: true, errorRates: []string{"0", "0"}, fakeMetric: fakeValStrs, satisfyStrs: satisfyStrs, fakeObj: fakeValStrs},
	{name: "experiment5", started: true, errorRates: []string{"0", "0"}, fakeMetric: fakeValStrs, satisfyStrs: satisfyStrs, fakeObj: fakeValStrs},
	{name: "experiment6", started: true, errorRates: []string{"0", "0"}, fakeMetric: fakeValStrs, satisfyStrs: satisfyStrs, fakeObj: fakeValStrs},
	{name: "experiment7", started: true, errorRates: []string{"0", "0"}, fakeMetric: fakeValStrs, satisfyStrs: satisfyStrs, fakeObj: fakeValStrs},
	{name: "experiment8", started: true, errorRates: []string{"0", "0"}, fakeMetric: fakeValStrs, satisfyStrs: satisfyStrs, fakeObj: fakeValStrs},
	{name: "experiment9", started: true, errorRates: []string{"0", "0"}, fakeMetric: fakeValStrs, satisfyStrs: satisfyStrs, fakeObj: fakeValStrs},
}

func init() {
	for i := 0; i < len(tests); i++ {
		e, err := getExp(tests[i].name)
		if err == nil {
			tests[i].exp = e
		} else {
			fmt.Println("Unable to extract experiment objects from files")
			os.Exit(1)
		}
	}
}

func TestExperiment(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// test Started()
			assert.Equal(t, tc.started, tc.exp.Started())
			// test GetVersions()
			if tc.exp.Started() {
				assert.Equal(t, []string{"default", "canary"}, tc.exp.GetVersions())
			} else {
				assert.Equal(t, []string([]string(nil)), tc.exp.GetVersions())
			}
			// test GetMetricValueStrs(...)
			assert.Equal(t, tc.errorRates, tc.exp.GetMetricValueStrs("error-rate"))
			assert.Equal(t, tc.fakeMetric, tc.exp.GetMetricValueStrs("fake-metric"))
			// test GetSatisfyStrs()
			assert.Equal(t, tc.satisfyStrs, tc.exp.GetSatisfyStrs(0))
			assert.Equal(t, tc.fakeObj, tc.exp.GetSatisfyStrs(10))
		})
	}
}

func TestGetMetricNameAndUnits(t *testing.T) {
	metricNameAndUnits := [4]string{"95th-percentile-tail-latency (milliseconds)", "mean-latency (milliseconds)", "error-rate", "request-count"}
	mnu := [4]string{}
	for i := 0; i < 4; i++ {
		mnu[i] = GetMetricNameAndUnits(tests[2].exp.Spec.Metrics[i])
	}
	assert.Equal(t, metricNameAndUnits, mnu)
}

func TestStringifyObjective(t *testing.T) {
	objectives := [2]string{"mean-latency <= 1000", "error-rate <= 0.010"}
	objs := [2]string{}
	for i := 0; i < 2; i++ {
		objs[i] = StringifyObjective(tests[2].exp.Spec.Criteria.Objectives[i])
	}
	assert.Equal(t, objectives, objs)
}
