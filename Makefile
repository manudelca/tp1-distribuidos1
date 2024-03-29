SHELL := /bin/bash
PWD := $(shell pwd)

docker-image:
	docker build -f ./metric-server/Dockerfile -t "metric-server:latest" .
	docker build -f ./test-client/Dockerfile -t "test-client:latest" .
	docker build -f ./test-client-metric/Dockerfile -t "test-client-metric:latest" .
	docker build -f ./test-client-query/Dockerfile -t "test-client-query:latest" .
	# Execute this command from time to time to clean up intermediate stages generated 
	# during client build (your hard drive will like this :) ). Don't left uncommented if you 
	# want to avoid rebuilding client image every time the docker-compose-up command 
	# is executed, even when client code has not changed
	# docker rmi `docker images --filter label=intermediateStageToBeDeleted=true -q`
.PHONY: docker-image

docker-compose-up: docker-image
	docker-compose -f docker-compose-dev.yaml up -d --build
.PHONY: docker-compose-up

docker-compose-down:
	docker-compose -f docker-compose-dev.yaml stop -t 90
	docker-compose -f docker-compose-dev.yaml down -t 90
.PHONY: docker-compose-down

docker-compose-logs:
	docker-compose -f docker-compose-dev.yaml logs -f
.PHONY: docker-compose-logs