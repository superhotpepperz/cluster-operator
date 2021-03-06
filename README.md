# StorageOS cluster-operator

[![Build Status](https://travis-ci.org/storageos/cluster-operator.svg?branch=master)](https://travis-ci.org/storageos/cluster-operator)
[![CircleCI](https://circleci.com/gh/storageos/cluster-operator.svg?style=svg)](https://circleci.com/gh/storageos/cluster-operator)

The StorageOS Cluster Operator deploys and configures a StorageOS cluster on
Kubernetes.

For quick installation of cluster operator, use the [cluster operator helm chart](https://github.com/storageos/charts/tree/master/stable/storageoscluster-operator).

## Pre-requisites

* Kubernetes 1.9+
* Kubernetes must be configured to allow (configured by default in 1.10+):
    - Privileged mode containers (enabled by default)
    - Feature gate: MountPropagation=true.  This can be done by appending `--feature-gates MountPropagation=true` to the
      kube-apiserver and kubelet services.

Refer to the [StorageOS prerequisites docs](https://docs.storageos.com/docs/prerequisites/overview) for more information.


## Setup/Development:

1. Install [operator-sdk](https://github.com/operator-framework/operator-sdk/tree/master#quick-start).
2. Run `operator-sdk generate k8s` if there's a change in api type.
3. Build operator container with `operator-sdk build storageos/cluster-operator:<tag>`
4. Apply the manifests in `deploy/` to install the operator
    - Apply `service_account.yaml` and `role_binding.yaml` to create a service account and to grant all the permissions.
    - Apply `crds/*_storageoscluster_crd.yaml` to define the custom resources.
    - Apply `operator.yaml` to install the operator. Change the container image in this file when installing a new operator.
    - Apply `crds/*_storageoscluster_cr.yaml` to create a `StorageOSCluster` custom resource.

**NOTE**: Installing StorageOS on Minikube is not currently supported due to missing [kernel prerequisites](https://docs.storageos.com/docs/prerequisites/systemconfiguration).

For development, run the operator outside of the k8s cluster by running:
```
$ make local-run
```

Build operator container image:
```
$ make image/cluster-operator OPERATOR_IMAGE=storageos/cluster-operator:test
```
This builds all the components and copies the binaries into the same container.

After creating a resource, query the resource:
```
$ kubectl get storageoscluster
NAME                READY     STATUS    AGE
example-storageos   3/3       Running   4m
```

## Inspect a StorageOSCluster Resource

Get all the details about the cluster:
```
$ kubectl describe storageoscluster/example-storageos
Name:         example-storageos
Namespace:    default
Labels:       <none>
Annotations:  kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"storageos.com/v1alpha1","kind":"StorageOSCluster","metadata":{"annotations":{},"name":"example-storageos","namespace":"default"},"spec":{"...
API Version:  storageos.com/v1alpha1
Kind:         StorageOSCluster
Metadata:
  Creation Timestamp:  2018-07-21T12:57:11Z
  Generation:          1
  Resource Version:    10939030
  Self Link:           /apis/storageos.com/v1alpha1/namespaces/default/storageosclusters/example-storageos
  UID:                 955b24a4-8ce5-11e8-956a-1866da35eee2
Spec:
  Join:  test07
Status:
  Node Health Status:
  ...
  ...
  Nodes:
    test09
    test08
    test07
  Phase:  Running
  Ready:  3/3
Events:   <none>
```

## StorageOSCluster Resource Configuration

Once the StorageOS operator is running, a StorageOS cluster can be deployed by
creating a Cluster Configuration. The parameters specified in the
configuration will define how StorageOS is deployed, the rest of the
installation details are handled by the operator.

The following tables lists the configurable spec parameters of the StorageOSCluster custom resource and their default values.

Parameter | Description | Default
--------- | ----------- | -------
`secretRefName` | Reference name of storageos secret |
`secretRefNamespace` | Namespace of storageos secret |
`namespace` | Namespace where storageos cluster resources are created | `storageos`
`images.nodeContainer` | StorageOS node container image | `storageos/node:1.0.0-rc4`
`images.initContainer` | StorageOS init container image | `storageos/init:0.1`
`images.csiDriverRegistrarContainer` | CSI Driver Registrar Container image | `quay.io/k8scsi/driver-registrar:v0.2.0`
`images.csiExternalProvisionerContainer` | CSI External Provisioner Container image | `quay.io/k8scsi/csi-provisioner:v0.3.0`
`images.csiExternalAttacherContainer` | CSI External Attacher Container image | `quay.io/k8scsi/csi-attacher:v0.3.0`
`csi.enable` | Enable CSI setup | `false`
`csi.enableProvisionCreds` | Enable CSI provision credentials | `false`
`csi.enableControllerPublishCreds` | Enable CSI controller publish credentials | `false`
`csi.enableNodePublishCreds` | Enable CSI node publish credentials | `false`
`service.name` | Name of the Service used by the cluster | `storageos`
`service.type` | Type of the Service used by the cluster | `ClusterIP`
`service.externalPort` | External port of the Service used by the cluster | `5705`
`service.internalPort` | Internal port of the Service used by the cluster | `5705`
`service.annotations` | Annotations of the Service used by the cluster |
`ingress.enable` | Enable ingress for the cluster | `false`
`ingress.hostname` | Hostname to be used in cluster ingress | `storageos.local`
`ingress.tls` | Enable TLS for the ingress | `false`
`ingress.annotations` | Annotations of the ingress used by the cluster |
`sharedDir` | Path to be shared with kubelet container when deployed as a pod | `/var/lib/kubelet/plugins/kubernetes.io~storageos`
`kvBackend.address` | Comma-separated list of addresses of external key-value store. (`1.2.3.4:2379,2.3.4.5:2379`) |
`kvBackend.backend` | Name of the key-value store to use. Set to `etcd` for external key-value store. | `embedded`
`pause` | Pause the operator for cluster maintenance | `false`
`debug` | Enable debug mode for all the cluster nodes | `false`
`nodeSelectorTerms` | Set node selector for storageos pod placement |
`resources` | Set resource requirements for the containers |


## Upgrading a StorageOS Cluster

An existing StorageOS cluster can be upgraded to a new version of StorageOS by
creating an Upgrade Configuration. The cluster-operator takes care of
downloading the new container image and updating all the nodes with new version
of StorageOS.
An example of of `StorageOSUpgrade` resource is [storageos_v1alpha1_storageosupgrade_cr.yaml](/deploy/crds/storageos_v1alpha1_storageosupgrade_cr.yaml).

Only offline upgrade is supported for now by cluster-operator. During the
upgrade, the applications that use StorageOS volumes are scaled down and the
whole StorageOS cluster is restarted with a new version. Once the StorageOS
cluster becomes usable, the applications are scaled up to their previous configuration.

Once an upgrade resource is created, events related to the upgrade can be
viewed in the upgrade object description. All the status and errors, if any,
encountered during the upgrade are posted as events.

```
$ kubectl describe storageosupgrades example-storageosupgrade
Name:         example-storageosupgrade
Namespace:    default
Labels:       <none>
Annotations:  kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"storageos.com/v1alpha1","kind":"StorageOSUpgrade","metadata":{"annotations":{},"name":"example-storageosupgrade","namespace":"default"},...
API Version:  storageos.com/v1alpha1
Kind:         StorageOSUpgrade
...
Spec:
  New Image:  storageos/node:1.0.0
Events:
  Type    Reason           Age   From                Message
  ----    ------           ----  ----                -------
  Normal  UpgradeComplete  21s   storageos-upgrader  StorageOS upgraded to storageos/node:1.0.0
```

## StorageOSUpgrade Resource Configuration

The following table lists the configurable spec parameters of the
StorageOSUpgrade custom resource and their default values.

Parameter | Description | Default
--------- | ----------- | -------
`newImage` | StorageOS node container image to upgrade to |


## TLS Support

To enable TLS, ensure that an ingress controller is installed in the cluster.
Set `ingress.enable` and `ingress.tls` to `true`.
Store the TLS cert and key as part of the storageos secret as:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: "storageos-api"
...
...
data:
  # echo -n '<secret>' | base64
  ...
  ...
  # Add base64 encoded TLS cert and key.
  tls.crt:
  tls.key:
```

## CSI

StorageOS also supports the [Container Storage Interface (CSI)](https://github.com/container-storage-interface/spec) to communicate with Kubernetes.
Only versions 1.10+ are supported. CSI ensures forward compatibility with future releases of
Kubernetes, as vendor-specific drivers will soon be deprecated from Kubernetes.
However, some functionalities are not supported yet.

To enable CSI, set `csi.enable` to `true` in the `StorageOSCluster` resource config.

```yaml
apiVersion: "storageos.com/v1alpha1"
kind: "StorageOSCluster"
metadata:
  name: "example-storageos"
  namespace: "default"
spec:
  secretRefName: "storageos-api"
  secretRefNamespace: "default"
  csi:
    enable: true
```

### CSI Credentials

To enable CSI Credentials, ensure that CSI is enabled by setting `csi.enable` to
`true`. Based on the type of credentials to enable, set the csi fields to `true`:
```yaml
apiVersion: "storageos.com/v1alpha1"
kind: "StorageOSCluster"
metadata:
  name: "example-storageos"
  namespace: "default"
spec:
  ...
  ...
  csi:
    enable: true
    enableProvisionCreds: true
    enableControllerPublishCreds: true
    enableNodePublishCreds: true
  ...
```

Specify the CSI credentials as part of the storageos secret object as:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: "storageos-api"
...
...
data:
  # echo -n '<secret>' | base64
  ...
  ...
  csiProvisionUsername:
  csiProvisionPassword:
  csiControllerPublishUsername:
  csiControllerPublishPassword:
  csiNodePublishUsername:
  csiNodePublishPassword:
```
