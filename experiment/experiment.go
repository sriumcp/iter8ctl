package experiment

import (
	"errors"
	"fmt"

	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
)

// Experiment is an enhancement of v2alpha1.Experiment struct. It provides various methods to build experiment description.
type Experiment struct {
	v2alpha1.Experiment
}

// Started indicates if at least one iteration of the experiment has completed
func (e *Experiment) Started() bool {
	return e.Status.CompletedIterations != nil && *e.Status.CompletedIterations > 0
}

// GetVersions returns the list of version names
func (e *Experiment) GetVersions() []string {
	versions := []string{e.Spec.VersionInfo.Baseline.Name}
	for _, c := range e.Spec.VersionInfo.Candidates {
		versions = append(versions, c.Name)
	}
	return versions
}

// RequestCountSpecified indicates if RequestCount metric is specified in the experiment
func (e *Experiment) RequestCountSpecified() bool {
	return e.Spec.Criteria.RequestCount != nil
}

// GetMetricStr returns the metric value (as a string) for a given metric and version
func (e *Experiment) GetMetricStr(metric string, version string) string {
	am := e.Status.Analysis.AggregatedMetrics
	if am == nil {
		return "unavailable"
	}
	if vals, ok := am.Data[metric]; ok {
		if val, ok := vals.Data[version]; ok {
			return fmt.Sprintf("%v", val.Value.ScaledValue(0))
		}
	}
	return "unavailable"
}

// GetRequestCountStrs returns the request count for each version
func (e *Experiment) GetRequestCountStrs() ([]string, error) {
	if !e.RequestCountSpecified() {
		return nil, errors.New("GetRequestCountStrs invoked on experiment without request count metric specification")
	}
	versions := e.GetVersions()
	reqs := make([]string, len(versions))
	for i, v := range versions {
		reqs[i] = e.GetMetricStr(*e.Spec.Criteria.RequestCount, v)
	}
	return reqs, nil
}
