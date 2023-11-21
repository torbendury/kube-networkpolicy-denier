.RECIPEPREFIX = >

# Variables
DEV_IMAGE_NAME := kube-networkpolicy-denier
RELEASE_IMAGE_NAME := torbendury/kube-networkpolicy-denier

# Create self signed certificates
ssl:
> hack/gen_certs.sh

# Run the server locally as a binary
run:
> go run cmd/main.go --cert=hack/webhook.pem --key=hack/webhook-key.pem

test:
> go test -v ./...

# Build the development stage container
dev:
> docker build --no-cache -t $(DEV_IMAGE_NAME):dev --target dev .
> docker run -it --rm -v ${PWD}:/app -p 8443:8443 $(DEV_IMAGE_NAME):dev

# Build the release container
build:
> docker build --no-cache -t $(RELEASE_IMAGE_NAME):latest --target release .

publish:
# TODO add versioning, maybe with GitHub action
# push to docker hub and release helm chart

release: test build publish

.PHONY: ssl run test dev release
