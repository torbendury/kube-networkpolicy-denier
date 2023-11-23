# 👮 kube-networkpolicy-denier

A very basic **admission controller** for Kubernetes that denies all network policies. It works as a validation webhook and can be used to prevent users from creating network policies. This is especially useful in multi-tenant clusters where you want to prevent users from creating network policies that might affect other users, or (like in my case) in environments where you exclusively rely on Istio AuthorizationPolicies.

## 🚀 Installation / Deployment

See the [README for the Helm Chart](helm/kube-networkpolicy-denier/README.md) for more information.

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
