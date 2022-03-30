DOCKER_IMAGE_NAME = protobuf-build-env
PROJECT_NAME = videocoin-network

REPOSITORY=gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME)

build:
	docker build -t $(REPOSITORY) .

push:
	docker push $(REPOSITORY)

clean:
	docker image rm $(REPOSITORY)
