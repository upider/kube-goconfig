TAG:=v0.0.1

all: build docker deploy

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build agent and manager binary.
	go build -o main cmd/main.go

docker-build: ## Build docker image
	cd docker && docker build -t 1445277435/kube-goconfig:$(TAG) -f Dockerfile ..

docker-push: ## Push docker image
	docker push 1445277435/kube-goconfig:$(TAG)

deploy: ## Deploy on kubernetes
	kubectl apply -f deployment/deployment.yaml

clean: ## Clean main and docker images
	rm -rf main
	docker rmi 1445277435/kube-goconfig:$(TAG) --force