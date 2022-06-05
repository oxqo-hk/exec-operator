# exec-operator

Execute "kubectl exec" commands on multiple pods selected by custom resource

## Description

- Execute shell command on selected pod.container[0].
- Pods are selected by IP, name and selector.
- All selecting factors will only bounded on CR’s namespace to prevent unintended actions.
- Same command will be executed on all pods selected by each factor
- Command will not be executed more than once even if a same pod selected multiple times by multiple factors.
- After attempt for execution, the command will marked as done, and not be executed again in reconciling loop, whether the command fails or not.
- The results of command can be found in CR’s Status

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Install Instances of Custom Resources:

```
kubectl apply -f config/samples/

```

1. Build and push your image to the location specified by `IMG`:

```
make docker-build docker-push IMG=<some-registry>/exec-operator:tag

```

1. Deploy the controller to the cluster with the image specified by `IMG`:

```
make deploy IMG=<some-registry>/exec-operator:tag

```

### Uninstall CRDs

To delete the CRDs from the cluster:

```
make uninstall

```

### Undeploy controller

UnDeploy the controller to the cluster:

```
make undeploy

```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster

### Test It Out

1. Install the CRDs into the cluster:

```
make install

```

1. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```
make run

```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```
make manifests

```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022 oxqo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

```
<http://www.apache.org/licenses/LICENSE-2.0>

```

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.