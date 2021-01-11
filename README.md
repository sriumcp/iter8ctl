[![Go Report Card](https://goreportcard.com/badge/github.com/iter8-tools/iter8ctl)](https://goreportcard.com/report/github.com/iter8-tools/iter8ctl)
[![Coverage](https://codecov.io/gh/iter8-tools/iter8ctl/branch/main/graphs/badge.svg?branch=main)](https://codecov.io/gh/iter8-tools/iter8ctl)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/iter8-tools/iter8ctl.svg)](https://pkg.go.dev/github.com/iter8-tools/iter8ctl)
# Iter8ctl
Iter8 command line utility for service operators to understand and diagnose their iter8 experiments.

Iter8ctl can be used with [iter8-kfserving](https://github.com/iter8-tools/iter8-kfserving) experiments.

# Installation
```
GOBIN=/usr/local/bin go install github.com/iter8-tools/iter8ctl
```
The above command installs `iter8ctl` under the `/usr/local/bin` directory. To install under a different directory, change the value of `GOBIN` above.

## Removal
```
rm <path-to-install-directory>/iter8ctl
```

# Usage

## Example 1
Describe an iter8 Experiment resource object present in your Kubernetes cluster.
```shell
kubectl get experiment sklearn-iris-experiment-1 -n kfserving-test -o yaml > experiment.yaml
iter8ctl describe -f experiment.yaml
```

## Example 2
Supply experiment YAML using console input.
```shell
kubectl get experiment sklearn-iris-experiment-1 -n kfserving-test -o yaml > experiment.yaml
cat experiment.yaml | iter8ctl describe -f -
```

## Example 3
Periodically fetch an iter8 Experiment resource object present in your Kubernetes cluster and describe it. You can change the frequency by adjusting the sleep interval below.
```shell
while clear; do
    kubectl get experiment sklearn-iris-experiment-1 -n kfserving-test -o yaml | iter8ctl describe -f -
    sleep 10.0
done
```

# Sample Output
The following is the output of executing `iter8ctl describe -f testdata/experiment8.yaml`; the `testdata` folder is part of the [iter8ctl GitHub repo](https://github.com/iter8-tools/iter8ctl) and contains sample experiments used in tests.

```shell
$ ./iter8ctl describe -f testdata/experiment8.yaml
******
Experiment name: sklearn-iris-experiment-1
Experiment namespace: kfserving-test
Experiment target: kfserving-test/sklearn-iris

******
Number of completed iterations: 10

******
Winning version: canary

******
Objectives
+--------------------------+---------+--------+
|        OBJECTIVE         | DEFAULT | CANARY |
+--------------------------+---------+--------+
| mean-latency <= 1000.000 | true    | true   |
+--------------------------+---------+--------+
| error-rate <= 0.010      | true    | true   |
+--------------------------+---------+--------+

******
Metrics
+--------------------------------+---------+---------+
|             METRIC             | DEFAULT | CANARY  |
+--------------------------------+---------+---------+
| 95th-percentile-tail-latency   | 330.682 | 310.320 |
| (milliseconds)                 |         |         |
+--------------------------------+---------+---------+
| mean-latency (milliseconds)    | 228.420 | 229.002 |
+--------------------------------+---------+---------+
| error-rate                     |   0.000 |   0.000 |
+--------------------------------+---------+---------+
| request-count                  | 117.445 |  57.715 |
+--------------------------------+---------+---------+
```

# Contributing

Documentation and code PRs are welcome. When contributing to this repository, please first discuss the change you wish to make using [Issues](https://github.com/iter8-tools/iter8ctl/issues), [Discussion](https://github.com/iter8-tools/iter8ctl/discussions), or [Slack](https://join.slack.com/t/iter8-tools/shared_invite/enQtODU0NTczMTQ5NDU4LTJmNGE1OTBhOWI4NzllZGE0ZjdhM2M3MzJlMjcxYjliMTJlM2YxMzQ4OWQ5NGViYTM2MTU4MWRkZTgxNzZiMzg).