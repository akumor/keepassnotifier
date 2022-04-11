.PHONY: docker_build
docker_build:
	IMAGE_NAME=$$REPOSITORY ./boilerplate/docker_build/docker_build.sh

.PHONY: dockerhub_push
dockerhub_push:
	IMAGE_NAME=akumor/$$REPOSITORY REGISTRY=docker.io ./boilerplate/docker_build/docker_build.sh