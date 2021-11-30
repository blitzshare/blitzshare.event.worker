install:
	go install golang.org/x/tools/cmd/goimports@latest
	go mod vendor

test:
	ENV=test && go test -v ./... -v -count=1

fix-format:
	gofmt -w -s app/ pkg/ cmd/ mocks/ testhelpers
	goimports -w app/ pkg/ cmd/ mocks/ testhelpers

start:
	go run cmd/main.go

build:
	GIN_MODE=release go build -o entrypoint cmd/main.go

k8s-apply:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment blitzshare-event-worker-deployment --namespace blitzshare-event-worker-ns

k8s-destroy:
	kubectl delete namespace blitzshare-event-worker-ns

build-deploy:
	make dockerhub-build
	make k8s-apply

docker-build:
	docker buildx build --platform linux/amd64 -t  blitzshare.api:latest .
	docker build -t  blitzshare.event.worker:latest .
	
dockerhub-build:
	make docker-build
	docker tag blitzshare.event.worker:latest iamkimchi/blitzshare.event.worker:latest
	docker push iamkimchi/blitzshare.event.worker:latest

minikube-svc:
	minikube service blitzshare-api-svc -n blitzshare-api-ns