# Overview

This repo contains example kustomize configs used to installed openshift features required for CNF workloads and a e2e functional test suite used to verify cnf related features.

## Contributing kustomize configs

All kustomize configs should be entirely declarative in nature. This means no bash plugin modules performing imparative tasks. Features should be installed simply by posting manifests to the cluster. After posting manifests, determining when the cluster has converged on those manifests successully should be observable.

## Usage

### Prerequisites

- You need a running OCP 4.4 cluster and a valid KUBECONFIG.
- You need at least one node with the `node-role.kubernetes.io/worker-cnf=""` label and a `MachineConfigPool` matching `worker-cnf` machine configurations

You can run `make setup-test-cluster` to have the first two (or the first in case of only one) workers labeled as `worker-cnf` and to have the `MachineConfigPool` created.

### Configuring

All the Makefile rules depend on two environment variable, either for deploying, waiting and choosing what tests to run.

##### FEATURES

i.e. `FEATURES="sctp ptp sriov"`, drives what features are going to be deployed using kustomize, and what tests are going to be run.

The current default values is `"sctp performace"`

##### FEATURES_ENVIRONMENT

i.e. `FEATURES_ENVIRONMENT=demo` determines the kustomization layer that will be used to deploy the choosen features.

The current default values is `e2e-gcp`

### Deployment

For each feature choosen via `FEATURES` we expect to have a layer either in [feature-configs/base](feature-configs/base) or in [feature-configs/$FEATURES_ENVIRONMENT](feature-configs/demo).

- run `FEATURES_ENVIRONMENT=demo make feature-deploy`.  
  This will try to apply all manifests in a loop until all deployments succeeded, or until it runs into a timeout.
- optionally run `FEATURES_ENVIRONMENT=demo make feature-wait` to be notified of when the features are deployed.

### Testing

We expect to have a section of [the test suite](functests/test_suite_test.go) named after each feature we want to test (for example [sctp](functests/sctp/sctp.go) named after the sctp feature).

In case a feature needs to be instrumented in a custom way, we expect a `test.sh` file under [/external-tests/$feature/test.sh](/external-tests/performance/test.sh).

A common use case is to rely on the feature's tier 1 tests taken from an external repo and [performance](/external-tests/performance/test.sh) provides a good example for that.

`FEATURES_ENVIRONMENT=demo make functests` is the command for executing the relevant tests.
