
SHELL := /bin/bash
CWD := $(shell cd -P -- '$(shell dirname -- "$0")' && pwd -P)
export GO111MODULE := on
export GOBIN := $(CWD)/.bin

install:
	go install $(shell go list -f '{{join .Imports " "}}' tools.go)
	go get -d github.com/vektra/mockery/v2/.../
	go mod vendor

test:
	go test  --tags='test' -v ./app/... -v -count=1 -cover -coverprofile=coverage.out
coverage-report-html:
	go tool cover -html=coverage.out		
fix-format:
	gofmt -w -s app/ cmd/ mocks/ testhelpers
	goimports -w app/  cmd/ mocks/ testhelpers
start:
	go run cmd/main.go
build:
	GIN_MODE=release go build -o entrypoint cmd/main.go
k8s-apply:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment blitzshare-event-worker-dpl --namespace blitzshare-ns
k8s-destroy:
	kubectl delete deployment blitzshare-event-worker-dpl
build-deploy:
	make dockerhub-build
	make k8s-apply
docker-build:
	docker buildx build --platform linux/amd64 -t  blitzshare.api:latest .
	docker build -t blitzshare.event.worker:latest .
	minikube image load blitzshare.event.worker:latest
	
dockerhub-build:
	make docker-build
	docker tag blitzshare.event.worker:latest iamkimchi/blitzshare.event.worker:latest
	docker push iamkimchi/blitzshare.event.worker:latest

minikube-svc:
	minikube service blitzshare-api-svc -n blitzshare-ns

build-mocks:
	.bin/mockery --all --dir "./app/"

.PHONY: test