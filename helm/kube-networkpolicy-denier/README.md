# kube-networkpolicy-denier

This chart installs a admission controller on a Kubernetes cluster using the Helm package manager.

It effectively denies all NetworkPolicies from being created on the cluster.

See changelog [here](https://github.com/torbendury/kube-networkpolicy-denier/blob/main/CHANGELOG.md).

## Source Code

You can find the source code of this chart on [GitHub](https://github.com/torbendury/kube-networkpolicy-denier).

## Installation

```bash
helm repo add kube-networkpolicy-denier https://torbendury.github.io/kube-networkpolicy-denier
helm repo update
helm install kube-networkpolicy-denier kube-networkpolicy-denier/kube-networkpolicy-denier --create-namespace --namespace kube-networkpolicy-denier
```

## Configuration

You can configure the following settings:

| Parameter                                       | Description                                                                                                                       | Default                                     |
|-------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------|
| `replicaCount`                                  | Number of replicas                                                                                                                | `2`                                         |
| `image.repository`                              | Image repository                                                                                                                  | `torbendury/kube-networkpolicy-denier`      |
| `image.tag`                                     | Image tag. Defaults to the Charts appVersion.                                                                                     | `""`                                        |
| `image.pullPolicy`                              | Image pull policy                                                                                                                 | `IfNotPresent`                              |
| `imagePullSecrets`                              | Image pull secrets                                                                                                                | `[]`                                        |
| `nameOverride`                                  | Override the name of the chart                                                                                                    | `""`                                        |
| `fullnameOverride`                              | Override the fullname of the chart                                                                                                | `""`                                        |
| `serviceAccount.create`                         | Specifies whether a ServiceAccount should be created                                                                              | `true`                                      |
| `serviceAccount.annotations`                    | Annotations to add to the ServiceAccount                                                                                          | `{}`                                        |
| `serviceAccount.name`                           | The name of the ServiceAccount. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template | `""`                                        |
| `serviceAccount.automount`                      | AutomountServiceAccountToken controls whether a service account token should be automatically mounted                             | `true`                                      |
| `podAnnotations`                                | Annotations to add to the pod                                                                                                     | `{}`                                        |
| `podSecurityContext`                            | Security context for the pod                                                                                                      | `{}`                                        |
| `podLabels`                                     | Labels to add to the pod                                                                                                          | `{}`                                        |
| `securityContext`                               | Security context for the container                                                                                                | `{}`                                        |
| `resources`                                     | Pod resource requests and limits                                                                                                  | `{}`                                        |
| `service.type`                                  | Kubernetes Service type                                                                                                           | `ClusterIP`                                 |
| `service.port`                                  | Kubernetes Service port                                                                                                           | `8443`                                      |
| `autoscaling.enabled`                           | Enable autoscaling for the deployment                                                                                             | `true`                                      |
| `autoscaling.minReplicas`                       | Minimum number of replicas                                                                                                        | `2`                                         |
| `autoscaling.maxReplicas`                       | Maximum number of replicas                                                                                                        | `5`                                         |
| `autoscaling.targetCPUUtilizationPercentage`    | Target CPU utilization percentage                                                                                                 | `80`                                        |
| `autoscaling.targetMemoryUtilizationPercentage` | Target memory utilization percentage                                                                                              | `80`                                        |
| `controller.response`                           | Response message to return when a NetworkPolicy is denied                                                                         | `"This webhook denies all NetworkPolicies"` |
