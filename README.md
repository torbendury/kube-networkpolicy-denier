[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/torbendury/kube-networkpolicy-denier.svg)](https://github.com/torbendury/kube-networkpolicy-denier)
[![Go Report Card](https://goreportcard.com/badge/github.com/torbendury/kube-networkpolicy-denier)](https://goreportcard.com/report/github.com/torbendury/kube-networkpolicy-denier)
![GitHub license](https://img.shields.io/github/license/torbendury/kube-networkpolicy-denier.svg)
![GitHub release](https://img.shields.io/github/release/torbendury/kube-networkpolicy-denier.svg)
[![GitHub latest commit](https://badgen.net/github/last-commit/torbendury/kube-networkpolicy-denier)](https://GitHub.com/torbendury/kube-networkpolicy-denier/commit/)

# 👮 kube-networkpolicy-denier

A very basic **admission controller** for Kubernetes that denies all network policies. It works as a validation webhook and can be used to prevent users from creating network policies. This is especially useful in multi-tenant clusters where you want to prevent users from creating network policies that might affect other users, or (like in my case) in environments where you exclusively rely on Istio AuthorizationPolicies.

See the full changelog [here](https://github.com/torbendury/kube-networkpolicy-denier/blob/main/CHANGELOG.md)

## 🚀 Installation / Deployment

See the [README for the Helm Chart](helm/kube-networkpolicy-denier/README.md) for more information.

## 🏎️ Resource Usage

This controller is very lightweight and does not consume a lot of resources. To give you a few key values, here are some numbers:

| CPU | Memory | Container Size | Load Test |
| ---------- | ------------- | -------------------- | --------- |
| 0.005 Cores | 5-10Mi | about 25MB | 1000 req/s delivered by single Pod: avg 1.83ms, min 149µs, max 700ms, p95 3.58ms |

### 🧑‍💻 Local Deployment

If you want to deploy `kube-networkpolicy-denier` locally, you can use the provided Makefile to spin up a Minikube cluster, build the container image yourself and deploy the controller to your cluster:

```bash
make local
```

## 📝 From Source

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### 💾 Prerequisites

Things you need to install the software from source:

- [Go](https://go.dev/doc/install) - The programming language used
- [Docker](https://docs.docker.com/get-docker/) - Containerization
- [Kubernetes](https://kubernetes.io/docs/tasks/tools/) - The container orchestration system used
  - This controller was tested on a Kubernetes 1.27.4 cluster, but should work on any Kubernetes cluster that supports admission webhooks in version `v1`. For instance, you can use a Minikube cluster for local development.
- [Helm](https://helm.sh/docs/intro/install/) (v3) - The package manager for Kubernetes, needed to deploy the chart.

### 🏃 Running

#### 💻 Locally

A step by step series of examples that tell you how to get a development environment running.

If you want to get started locally fast, you can use the Makefile to generate self-signed certificates and start the server locally:

```bash
make ssl
make run
```

If you want to work in a containerized environment to prevent pollution of your machine, you can use the provided Dockerfile to build a container image and run it locally:

```bash
make dev
```

From this running container on, you can proceed and build everything.

## 🤝 Contributing

Contributions are very welcome. I am happy to accept pull requests or issues. Please stay respectful. I don't plan on adding a code of conduct, but please be nice.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. Basically, you can do whatever you want with this project, but you have to include the license and the license notice. And if you break something while using this piece of software or anything around it, I'm not responsible for that.

## 🫴 Versioning

This project uses a mixed approach for versioning.

### 📦 Helm Chart

The version of the **Helm Chart** uses [Semantic Versioning](https://semver.org).

### 🐳 Docker Image

The version of the **controller Docker image** uses the **short git commit hash** of the commit that was used to build the respective image. A newer build of the image will always have a newer commit hash.

### Mixed Versioning

A newer build is not guaranteed to make it into a new release of the Helm Chart. This is mostly due to the fact that not every change to the controller is relevant for a release of the Helm Chart or not releasable at all. Anyway, you are free to use any version of the controller image with any version of the Helm Chart. Although, it is recommended to use the controller image version which is referenced in the `appVersion` field of the Helm Chart. A image version that has not been referenced in a Helm Chart version is not guaranteed to work as expected.

### 🛳️ Kubernetes Versions

I tested this Helm Chart in multiple Kubernetes *v1.27* clusters. According to the [Kubernetes API Version Guide](https://kubernetes.io/docs/reference/using-api/deprecation-guide/), the API that was introduced the last was `autoscaling/v2` since Kubernetes *v1.23*. This means that this Helm Chart should work with Kubernetes *v1.23* and above and can be recommended.

If you want to use this Helm Chart with a Kubernetes version below *v1.23*, you will need to disable the `HorizontalPodAutoscaler` by setting the `autoscaling.enabled` field in the `values.yaml` file to `false`. If you do so, you *should* be able to run the Helm Chart with any Kubernetes version above (and including) *v1.16*. However, this is not recommended because I did **NOT** test with Kubernetes *v1.16* and I can't guarantee that everything will work as expected in the long run. (*v1.16* was EOL in 2020, please upgrade your cluster if you are still using it.)
