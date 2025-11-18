DOCKERHUB_USER ?= tuilatung2001
APP_NAME       ?= todo-list-be

DEV_VERSION    ?= 0.0.1
PROD_VERSION   ?= 0.1.0

IMAGE          ?= $(DOCKERHUB_USER)/$(APP_NAME):$(DEV_VERSION)
PROD_IMAGE     ?= $(DOCKERHUB_USER)/$(APP_NAME):$(PROD_VERSION)
PROD_LATEST    ?= $(DOCKERHUB_USER)/$(APP_NAME):latest

.PHONY: build test push build-prod push-prod

build: ## Build docker image (dev)
	@echo "Building DEV image: $(IMAGE)"
	docker build -t $(IMAGE) .

test: ## Run Go tests locally
	go test ./...

push: ## Push dev images
	@echo "Pushing DEV image..."
	docker push $(IMAGE)

build-prod: ## Build production image
	@echo "Building PROD image: $(PROD_IMAGE)"
	docker build -t $(PROD_IMAGE) .
	docker tag $(PROD_IMAGE) $(PROD_LATEST)

push-prod: ## Push production images
	@echo "Pushing PROD images..."
	docker push $(PROD_IMAGE)
	docker push $(PROD_LATEST)
