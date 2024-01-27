.RECIPEPREFIX = >
.PHONY: ssl run test build kube helm local stress

### Variables
DEV_IMAGE_NAME := kube-networkpolicy-denier
RELEASE_IMAGE_NAME := torbendury/kube-networkpolicy-denier

### Create self signed certificates
ssl:
> hack/gen_certs.sh

### Run the server locally as a binary
run:
> cd cmd; go run main.go http.go admission.go --cert=../hack/webhook.pem --key=../hack/webhook-key.pem

### Run the tests
test:
> go test -v -race ./...
> helm lint helm/kube-networkpolicy-denier
> go mod verify
> go vet ./...
> go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
> go run golang.org/x/vuln/cmd/govulncheck@latest ./...

### Run performance benchmarks
bench:
> go test ./... -bench . -count=5 -benchmem | grep -v INFO

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

miniload:
> minikube image load $(RELEASE_IMAGE_NAME):latest

### Install the Helm Chart
helm:
> helm upgrade --install kube-networkpolicy-denier ./helm/kube-networkpolicy-denier --set image.tag=latest --namespace kube-networkpolicy-denier --create-namespace

github: build miniload helm
> sleep 5
> kubectl diff -f hack/netpol.yml; if [ $$? -ne 2 ]; then echo "NetworkPolicy was not denied correctly"; false; fi

### Run a series of stress tests
stress:
> k6 run hack/k6/smoke.js
> sleep 10
> k6 run hack/k6/soak.js
> sleep 10
> k6 run hack/k6/stress.js
> sleep 10
> k6 run hack/k6/spike.js
> sleep 10
> k6 run hack/k6/breakpoint.js

### Create CPU profile and show pprof http server
profile:
> go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof -benchmem ./...
> go tool pprof -http=:8080 cpu.prof

### Create a complete local environment
local: build kube helm
