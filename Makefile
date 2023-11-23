.RECIPEPREFIX = >
.PHONY: ssl run test dev build kube helm local

### Variables
DEV_IMAGE_NAME := kube-networkpolicy-denier
RELEASE_IMAGE_NAME := torbendury/kube-networkpolicy-denier

### Create self signed certificates
ssl:
> hack/gen_certs.sh

### Run the server locally as a binary
run:
> go run cmd/main.go --cert=hack/webhook.pem --key=hack/webhook-key.pem

### Run the tests
test:
> go test -v ./...
> helm lint helm/kube-networkpolicy-denier

### Build the development stage container
dev:
> docker build --no-cache -t $(DEV_IMAGE_NAME):dev --target dev .
> docker run -it --rm -v ${PWD}:/app -p 8443:8443 $(DEV_IMAGE_NAME):dev

### Build the release container
# This is not used in the CI/CD pipeline!
build:
> docker build --no-cache -t $(RELEASE_IMAGE_NAME):latest --target release .

### Create local test cluster
# This is not used in the CI/CD pipeline!
kube:
> minikube start
> minikube image load $(RELEASE_IMAGE_NAME):latest
> sleep 10

### Install the Helm Chart
helm:
> helm install kube-networkpolicy-denier kube-networkpolicy-denier/kube-networkpolicy-denier --set image.tag=latest --version 0.0.3 --namespace kube-networkpolicy-denier --create-namespace

### Create a complete local environment
local: build kube helm