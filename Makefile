REGISTRY=dvision
IMAGE_NAME=btckey-gen
VERSION=0.0.1

build:
	go build -o ./bin/$(IMAGE_NAME) ./cmd/

run: build
	X_API_KEY='' ./bin/$(IMAGE_NAME)


image:
	docker build -t $(REGISTRY)/${IMAGE_NAME}:$(VERSION) -f ./Dockerfile .

