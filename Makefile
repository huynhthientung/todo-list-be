DOCKERHUB_USER ?= toilatung2001
APP_NAME       ?= todo-list-be
BUILD_NUMBER   ?= 0.0.1

IMAGE          ?= $(DOCKERHUB_USER)/$(APP_NAME):$(BUILD_NUMBER)
LATEST_IMAGE   ?= $(DOCKERHUB_USER)/$(APP_NAME):latest

CHART_PATH     ?= helm/$(APP_NAME)

.PHONY: build push deploy

## Build docker image
build:
	@echo "Building: $(IMAGE)"
	docker build -t $(IMAGE) .
	docker tag $(IMAGE) $(LATEST_IMAGE)

## Push Docker Hub
push:
	@echo "Pushing to Docker Hub..."
	docker push $(IMAGE)
	docker push $(LATEST_IMAGE)

## Deploy with Helm
deploy:
	@echo "Deploying using Helm..."
	helm upgrade --install $(APP_NAME) $(CHART_PATH) \
		--set image.repository=$(DOCKERHUB_USER)/$(APP_NAME) \
		--set image.tag=$(BUILD_NUMBER) \
