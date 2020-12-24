package exp

import (
	v2alpha1 "github.com/iter8-tools/etc3/api/v2alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Experiment corresponds to iter8 Experiment type
type Experiment v2alpha1.Experiment

// GetExperiment gets the Kubernetes experiment object.
func GetExperiment(name *string, namespace *string, apiVersion *string) (Experiment, error) {
	return Experiment{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{},
		Spec:       v2alpha1.ExperimentSpec{},
		Status:     v2alpha1.ExperimentStatus{},
	}, nil
}

// PrintAnalysis does nothing for now.
func (e *Experiment) PrintAnalysis() {}
