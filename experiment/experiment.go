package experiment

import (
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

// GetMetricStr returns the metric value (as a string) for a given metric and version
func (e *Experiment) GetMetricStr(metric string, version string) string {
	am := e.Status.Analysis.AggregatedMetrics
	if am == nil {
		return "unavailable"
	}
	if vals, ok := am.Data[metric]; ok {
		if val, ok := vals.Data[version]; ok {
			return val.Value.AsDec().String()
		}
	}
	return "unavailable"
}

// GetMetricValueStrs returns the given metric's value for each version
func (e *Experiment) GetMetricValueStrs(metric string) []string {
	versions := e.GetVersions()
	reqs := make([]string, len(versions))
	for i, v := range versions {
		reqs[i] = e.GetMetricStr(metric, v)
	}
	return reqs
}

// GetMetricNameAndUnits from metric info
func GetMetricNameAndUnits(metricInfo v2alpha1.MetricInfo) string {
	r := metricInfo.Name
	if metricInfo.MetricObj.Spec.Units != nil {
		r += fmt.Sprintf(" (" + *metricInfo.MetricObj.Spec.Units + ")")
	}
	return r
}

// StringifyObjective returns a string representation of objective (with <= notation)
func StringifyObjective(objective v2alpha1.Objective) string {
	r := ""
	if objective.LowerLimit != nil {
		r += objective.LowerLimit.AsDec().String() + " <= "
	}
	r += objective.Metric
	if objective.UpperLimit != nil {
		r += " <= " + objective.UpperLimit.AsDec().String()
	}
	return r
}

// GetSatisfyStr returns a true/false/unavailable valued string denotating if a version satisfies the objective
func (e *Experiment) GetSatisfyStr(objectiveIndex int, version string) string {
	va := e.Status.Analysis.VersionAssessments
	if va == nil {
		return "unavailable"
	}
	if vals, ok := va.Data[version]; ok {
		if len(vals) > objectiveIndex {
			return fmt.Sprintf("%v", vals[objectiveIndex])
		}
	}
	return "unavailable"
}

// GetSatisfyStrs returns a slice of true/false/unavailable valued strings for an objective denoting if it is satisfied by versions
func (e *Experiment) GetSatisfyStrs(objectiveIndex int) []string {
	versions := e.GetVersions()
	sat := make([]string, len(versions))
	for i, v := range versions {
		sat[i] = e.GetSatisfyStr(objectiveIndex, v)
	}
	return sat
}
