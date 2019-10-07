Developers
=============
This document describe how to build and run this project locally.

There are two parts to this project:
* The Kubernetes Operator
* The endpoint-runner

# The respository structure
At your `$GOPATH` you should create a directory named `managedkube.com`

In the `managedkube.com` directory, checkout this repository into there.

# The Kubernetes Operator
The Kubernetes Operator handles gathering the info for the `endpoint-runner` and passing it the config on what
to test.  This operator watches the configured `ingresses` on the Kubernetes system and launches the `endpoint-runner`
Kubernetes Deployment with the correct configuration.  When Ingresses changes via (add/update/delete) this Operator
will be notified and it will adjust the `endpoint-runner` accordingly.

The Kubernetes Operator was created with the `operator-framework/operator-sdk`.

Instruction on the creation is located in the file: [docs/operator-usage.md](./operator-usage.md)

You should read this doc over to get a feel on how it was created.

## Running the Operator

1. Deploy out the CRD

You first have to deploy out the CRD to the Kubernetes cluster.

2. Run the Operator

Instructions can be found here: [docs/operator-usage.md#deploy-the-crd-to-the-kubernetes-cluster](./operator-usage.md#deploy-the-crd-to-the-kubernetes-cluster)
  
# The endpoint-runner
The `endpoint-runner` is the process that reaches out to the HTTP endpoint(s) to check if it is alive.  

For input it reads an envar (ENDPOINT_TEST_JSON) for configuration on what it should check.

## Building

pwd: `./cmd/endpoint-runner`

Run:
```
go build
```

## Running

```
export ENDPOINT_TEST_JSON={}
./endpoint-runner
```

