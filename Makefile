
REGISTRY:=npodewitz
IMAGE_NAME:=e3w
CONTAINER_NAME:=${IMAGE_NAME}
DOCKER_RUN_ARGS:=-p 8081:8080
VERSION:=


.PHONY: docker-build docker-build-nc run debug debug-exec stop up up-debug clean build

all: docker-build

dep-build: dep build

build:
	$(MAKE) -C src/ $@
#	export GOPATH="~/go/"
#	CGO_ENABLED=0 go build src/

dep:
	$(MAKE) -C src/ $@
#	export GOPATH="~/go/"
#	rm -rf src/vendor/*
#	dep ensure

docker-build:
	docker build -t ${IMAGE_NAME} .

docker-build-nc:
	docker build --no-cache -t ${IMAGE_NAME} .

run:
	docker run --name ${CONTAINER_NAME} ${DOCKER_RUN_ARGS} ${IMAGE_NAME}

debug:
	docker run -it --name ${CONTAINER_NAME} --entrypoint /bin/bash ${IMAGE_NAME}

debug-exec:
	docker exec -it ${CONTAINER_NAME} /bin/bash

stop:
	-docker stop ${CONTAINER_NAME}

up: clean build run

up-debug: clean build-debug run

clean: stop
	-docker rm -v ${CONTAINER_NAME}
