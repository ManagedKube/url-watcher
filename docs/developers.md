Developers
=============
This document describe how to build and run this project locally.

There are two parts to this project:
* The Kubernetes Operator
* The endpoint-runner

# The respository structure
At your $GOPATh you should create a directory named `managedkube.com`

In the `managedkube.com` directory, checkout this repository into there.

# The Kubernetes Operator
The Kubernetes Operator was created with the `operator-framework/operator-sdk`.

Instruction on the creation is located in the file: [docs/operator-usage.md](./operator-usage.md)

You should read this doc over to get a feel on how it was created.

## Running the Operator

1. Deploy out the CRD

You first have to deploy out the CRD to the Kubernetes cluster.

2. Run the Operator

Instructions can be found here: [docs/operator-usage.md#deploy-the-crd-to-the-kubernetes-cluster](./operator-usage.md#deploy-the-crd-to-the-kubernetes-cluster)
  
