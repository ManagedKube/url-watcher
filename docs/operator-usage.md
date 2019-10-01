# Operator Usage

## The initial creation:
Following this guide: https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#2-run-locally-outside-the-cluster


### Create the project scaffolding

```
operator-sdk new url-watcher
```

Since we already had the repo, I created this and then copied back over the `.git` directory back over to the
newly created repo.


### Create a new CRD

```
operator-sdk add api --api-version=url-watcher.managedkube.com/v1alpha1 --kind=Watcher
```

This will create `./pkg/apis/urlwatcher/v1alpha1/watcher_types.go`

After modifying the `*_types.go` file always run the following command to update the generated code for that resource type:

```
$ operator-sdk generate k8s
```

### Add a new controller

```
operator-sdk add controller --api-version=url-watcher.managedkube.com/v1alpha1 --kind=Watcher
```

Create a new controller file `./pkg/controller/watcher/watcher_controller.go`

### Deploy the CRD to the Kubernetes cluster

The CRD needs to be created first before running the operator or it will fail

```
kubectl apply -f deploy/crds/urlwatcher_v1alpha1_watcher_crd.yaml
```

### Running the controller

```
export OPERATOR_NAME=url-watcher
operator-sdk up local --namespace=default
```