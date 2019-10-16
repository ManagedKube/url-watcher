url-watcher
============

# Development Setup

## Main getting started guide

https://github.com/operator-framework/operator-sdk/blob/master/doc/helm/user-guide.md

## operator-sdk

Downloading the operator-sdk binary:

https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md

## Operator Types

### Helm Operator

https://github.com/operator-framework/operator-sdk/blob/master/doc/helm/user-guide.md

This walks you through how to wrap the operator around a Helm chart.  The Helm chart can be a custom chart you
create or an existing chart locally or from a Helm repository like helm/stable.

It then allows you to control all of the Helm values via CRDs.

### Custom Operator

https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md

Need Go version: 1.12


```
export GO111MODULE=on
```


```
export GOROOT=/usr/local/go-1.12

export PATH=/usr/local/go-1.12/bin:/home/g44/google-cloud-sdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/usr/local/go-1.12/bin:/home/g44/go/bin
```

Set kubeconfig
```
export KUBECONFIG=~/.kube/config
```


Run:
```
operator-sdk up local --namespace=default
```

# Using the debugger

https://github.com/operator-framework/operator-sdk/issues/1315


# Go modules

https://blog.golang.org/using-go-modules

https://medium.com/@fonseka.live/getting-started-with-go-modules-b3dac652066d

## Setting up the project

```
go mod init
```

# endpoint-runner

Running
```
cd endpoint-runner
go run main.go
```