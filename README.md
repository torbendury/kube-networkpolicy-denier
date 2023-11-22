# kube-networkpolicy-denier

A very basic **admission controller** for Kubernetes that denies all network policies. It works as a validation webhook and can be used to prevent users from creating network policies. This is especially useful in multi-tenant clusters where you want to prevent users from creating network policies that might affect other users, or (like in my case) in environments where you exclusively rely on Istio AuthorizationPolicies.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

You have multiple possibilities of running this project:

- Build the binary and the container yourself and deploy it to your cluster
- Use the pre-built container image from DockerHub

### Prerequisites

Things you need to install the software and how to install them:

- [Go](https://go.dev/doc/install) - The programming language used, only needed if you want to build the binary yourself
- [Docker](https://docs.docker.com/get-docker/) - Containerization, only needed if you want to build the container yourself
- [Kubernetes](https://kubernetes.io/docs/tasks/tools/) - The container orchestration system used
  - This controller was tested on a Kubernetes 1.27.4 cluster, but should work on any Kubernetes cluster that supports admission webhooks in version `v1`.
- [Helm](https://helm.sh/docs/intro/install/) (v3) - The package manager for Kubernetes, needed to deploy the chart

### Installing

#### Installing via Helm

See the [README for the Helm Chart](helm/kube-networkpolicy-denier/README.md) for more information.

#### Locally

A step by step series of examples that tell you how to get a development environment running.

If you want to get started locally fast, you can use the Makefile to generate self-signed certificates and start the server locally:

```bash
make ssl
make run
```

If you want to work in a containerized environment to prevent pollution of your local environment, you can use the provided Dockerfile to build a container image and run it locally:

```bash
make dev
```

If you want to build the container image yourself, you can use the following command:

**Work in Progress!**

```bash
make release
```

## Deployment

See the [README for the Helm Chart](helm/kube-networkpolicy-denier/README.md) for more information.

## Contributing

Contributions are very welcome. I am happy to accept pull requests or issues. Please stay respectful. I don't plan on adding a code of conduct, but please be nice.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. Basically, you can do whatever you want with this project, but you have to include the license and the license notice. And if you break something, it's not my fault.
