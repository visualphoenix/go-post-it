SHELL := /bin/bash
IMAGE_NAME ?= app
DOCKER_USERNAME ?= visualphoenix

.PHONY: all clean build builder

build: Dockerfile build/$(IMAGE_NAME)-amd64 build/upload.gtpl
	docker build \
		--build-arg app=$(IMAGE_NAME) \
		-f Dockerfile \
		-t $(IMAGE_NAME) .

builder: Dockerfile.builder
	docker build \
		--build-arg app=$(IMAGE_NAME) \
		--build-arg user=$(DOCKER_USERNAME) \
		-f Dockerfile.builder \
		-t $(IMAGE_NAME)-builder .

build/upload.gtpl:
	mkdir -p $(@D) && \
	cp $(shell basename $(@)) $(@) && \
	touch $(@)

build/$(IMAGE_NAME)-amd64: builder
	set -x && \
	mkdir -p $(@D) && \
	pushd $(@D) &>/dev/null && \
	docker run \
		--rm $(IMAGE_NAME)-builder \
		sh -c 'cd /go/bin && tar -c $(IMAGE_NAME)*' | tar -xv && \
	popd &>/dev/null && \
	touch $(@)

clean:
	test -d build && rm -r build || true
