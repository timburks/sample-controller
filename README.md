# api-controller

This repository implements a simple controller for watching a collection of
API-related resources that are defined with CustomResourceDefinitions (CRDs).

It is directly based on the [sample-controller](https://github.com/kubernetes/sample-controller) example.

It makes use of the generators in [k8s.io/code-generator](https://github.com/kubernetes/code-generator)
to generate a typed client, informers, listers and deep-copy functions. You can
do this yourself using the `./hack/update-codegen.sh` script.

The `update-codegen` script will automatically generate the following files &
directories:

* `pkg/apis/apicontroller/v1alpha1/zz_generated.deepcopy.go`
* `pkg/generated/`

Changes should not be made to these files manually, and when creating your own
controller based off of this implementation you should not copy these files and
instead run the `update-codegen` script to generate your own.

## Details

The controller uses [client-go library](https://github.com/kubernetes/client-go/tree/master/tools/cache) extensively.
The details of interaction points of the controller with various mechanisms from this library are
explained [here](docs/controller-client-go.md).

## Fetch api-controller and its dependencies

When using go 1.11 modules (`GO111MODULE=on`), issue the following
commands --- starting in whatever working directory you like.

```sh
git clone https://github.com/kubernetes/api-controller.git
cd api-controller
```

Note, however, that if you intend to
generate code then you will also need the
code-generator repo to exist in an old-style location.  One easy way
to do this is to use the command `go mod vendor` to create and
populate the `vendor` directory.

## Purpose

This is an exploration of how to build a kube-like controller with multiple types.

## Running

**Prerequisite**: Since the api-controller uses `apps/v1` deployments, the Kubernetes cluster version should be greater than 1.9.

```sh
# assumes you have a working kubeconfig, not required if operating in-cluster
go build -o api-controller .
./api-controller -kubeconfig=$HOME/.kube/config

# create the CustomResourceDefinitions
kubectl apply -f artifacts/examples/crds -R

# create the custom resources
kubectl apply -f artifacts/examples/apis -R

# check api products 
kubectl get apiproducts
```

## Subresources

Custom Resources support `/status` and `/scale` [subresources](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#subresources). The `CustomResourceSubresources` feature is in GA from v1.16.

### Example

The CRDs in [`crds`](./artifacts/examples/crds) enable the `/status` subresource for custom resources.
This means that [`UpdateStatus`](./controller.go) can be used by the controller to update only the status part of the custom resources.

To understand why only the status part of the custom resource should be updated, please refer to the [Kubernetes API conventions](https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status).


## A Note on the API version
The [group](https://kubernetes.io/docs/reference/using-api/#api-groups) version of the custom resources is `v1alpha`, this can be evolved to a stable API version, `v1`, using [CRD Versioning](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/).

## Cleanup

You can clean up the created CustomResourceDefinitions with:
```sh
kubectl delete crd apiproducts.apigee.dev
kubectl delete crd apiversions.apigee.dev
kubectl delete crd apidescriptions.apigee.dev
kubectl delete crd apideployments.apigee.dev
kubectl delete crd apiartifacts.apigee.dev
```