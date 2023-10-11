# Enable and disable DR

*Ramen* integrates with *OCM* to protect workloads on managed clusters.
To enable DR for an application, add a `DRPlacementControl` resource,
specifying how *Ramen* protects the application. Removing the
`DRPlacementControl` disables DR for the application.

## Enable DR for an application

When DR is enabled for an application *Ramen* may failover or relocate
the application to another managed cluster per user request. To make
this possible, *OCM* scheduling must be disabled for the application.

To disable *OCM* scheduling, add the following annotation to the
application placement:

```yaml
  annotations:
    cluster.open-cluster-management.io/experimental-scheduling-disable: "true"
```

At this point *OCM* will not change the application placement, and we
can add the `DRPlacementControl` resource.

To create the `DRPlacementControl` you need to look up some details
about the application:

- The name of managed cluster the application is running on, selected by
  *OCM* when the application was deployed (example: `dr1`)
- The application namespace where the `Placement` resource is deployed
  (example: `busybox-sample`)
- The name of the `DRPolicy` to use, created when configuring *Ramen*
  (example: `dr-policy`)
- The name of the `Placement` resource (example: `busybox-placement`)
- The label selector for selecting application `pvcs` for replication
  (example: `appname=busybox`)

Here is a sample drpc resource:

```yaml
apiVersion: ramendr.openshift.io/v1alpha1
kind: DRPlacementControl
metadata:
  name: busybox-drpc
  namespace: busybox-sample
  labels:
    app: busybox-sample
spec:
  preferredCluster: cluster1
  drPolicyRef:
    name: dr-policy
  placementRef:
    kind: Placement
    name: busybox-placement
  pvcSelector:
    matchLabels:
      appname: busybox
```

Applying this resource will enable DR for the application.
