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
	kubectl apply -f k8s/config/namespace.yaml
	kubectl apply -f k8s/config/deployment.yaml
	kubectl apply -f k8s/config/service.yaml
	kubectl apply -f k8s/config/redis.yml
	kubectl rollout restart deployment blitzshare-api-deployment --namespace blitzshare-api-ns

k8s-destroy:
	kubectl delete namespace blitzshare-api-ns


build-deploy:
	make dockerhub-build
	make k8s-apply


dockerhub-build:
	# docker buildx build --platform linux/amd64 -t  blitzshare.api:latest .
	docker build -t  blitzshare.api:latest .
	docker tag blitzshare.api:latest iamkimchi/blitzshare.api:latest
	docker push iamkimchi/blitzshare.api:latest

minikube-svc:
	minikube service blitzshare-api-svc -n blitzshare-api-ns