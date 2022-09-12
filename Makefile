IMAGE_NAME := "helloworld"

.PHONY: build
build:
	go build -o bin/helloworld

.PHONY: run
run:
	@go run main.go

.PHONY: test
test:
	@go test ./...

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: docker-build
docker-build:
	@docker build -f Dockerfile -t $(IMAGE_NAME) .

.PHONY: docker-run
docker-run:
	@docker run -p 127.0.0.1:8080:8080/tcp $(IMAGE_NAME)

.PHONY: docker-up
docker-up: docker-build docker-run

.PHONY: minikube-setup
minikube-setup:
	@which minikube || (echo Please install Minikube; exit 1)
	@minikube ip 2>&1 || ( \
		echo "Starting up Minikube"; \
		minikube start \
	)
	@minikube image load $(IMAGE_NAME)

.PHONY: deploy
deploy:
	@helm upgrade --install helloworld deploy/helloworld --values deploy/helloworld/values.yaml

.PHONY: deploy-test
deploy-test:
	@helm test helloworld
