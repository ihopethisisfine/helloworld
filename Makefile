IMAGE_NAME := "helloworld"
IMAGE_TAG  := "latest"

help: ## Print this helpful message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-21s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Builds golang binary for the helloworld API
	go build -o bin/helloworld

.PHONY: run
run: ## Runs the API locally
	@go run main.go

.PHONY: test
test: ## Runs all tests for the API
	@go test ./...

.PHONY: clean
clean: ## Removes the binary
	@rm -rf bin

.PHONY: docker-build
docker-build: ## Builds helloworld docker image with tag latest by default
	@docker build -f Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: docker-run
docker-run: ## Runs helloworld docker image mapping port 8080
	@docker run -p 127.0.0.1:8080:8080/tcp $(IMAGE_NAME)

.PHONY: docker-up
docker-up: docker-build docker-run ## Builds and runs helloworld docker image mapping port 8080

.PHONY: docker-push
docker-push: ## Pushes helloworld docker image
	@docker push $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: minikube-setup
minikube-setup: ## Sets up minikube
	@which minikube || (echo Please install Minikube; exit 1)
	@minikube ip 2>&1 || ( \
		echo "Starting up Minikube"; \
		minikube start \
	)

.PHONY: deploy-chart
deploy-chart: ## Deploys latest helm chart published version using localDynamoDB by default
	@helm repo add helloworld https://ihopethisisfine.github.io/helloworld
	@helm repo update helloworld
	@helm upgrade --install helloworld helloworld/helloworld --set localDynamodb.enable=true

.PHONY: deploy-prod-chart
deploy-prod-chart: ## Deploys latest helm chart published with the latest published docker image (assumes it is running on CI)
	@helm repo add helloworld https://ihopethisisfine.github.io/helloworld
	@helm repo update helloworld
	@helm upgrade --install helloworld helloworld/helloworld --set image.tag=${GITHUB_SHA} --set serviceAccount.create=true --set serviceAccount.name=helloworld --set ingress.enabled=true

.PHONY: deploy-local-chart
deploy-local-chart: ## Deploys helm chart with local changes using localDynamoDB by default
	@helm upgrade --install helloworld charts/helloworld --values charts/helloworld/values.yaml --set localDynamodb.enable=true

.PHONY: deploy-test
deploy-test: ## Runs a connection test with helm
	@helm test helloworld

.PHONY: compose-up
compose-up: ## Launches helloworld api and local dynamo using docker compose
	@docker-compose up -d

.PHONY: compose-down
compose-down: ## Tears down helloworld with docker compose
	@docker-compose down
