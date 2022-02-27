
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
acceptance-tests:
	cd "$(CWD)/test"  &&  ../.bin/godog ./**/*.feature
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
build-mocks:
	.bin/mockery --all --dir "./app/"
.PHONY: test