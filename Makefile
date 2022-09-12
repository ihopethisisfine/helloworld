IMAGE_NAME := "helloworld"
IMAGE_TAG  := "latest"

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
	@docker build -f Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: docker-run
docker-run:
	@docker run -p 127.0.0.1:8080:8080/tcp $(IMAGE_NAME)

.PHONY: docker-up
docker-up: docker-build docker-run

.PHONY: docker-push
docker-push:
	@docker push $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: deploy
deploy:
	@helm upgrade --install helloworld charts/helloworld --values deploy/helloworld/values.yaml

.PHONY: deploy-test
deploy-test:
	@helm test helloworld

.PHONY: compose-up
compose-up:
	@docker-compose up -d

.PHONY: compose-down
compose-down:
	@docker-compose down
