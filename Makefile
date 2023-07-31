DOCKER_CONFIG_PATH := docker/docker-compose.yaml
BUILD_CONTEXT := .
GO_CONFIG := ./config/local.yaml

# Build the Docker container
docker-build:
	sudo docker build -t clickhouse -f $(DOCKERFILE_PATH) $(BUILD_CONTEXT)

docker-run:
	sudo docker compose -f $(DOCKER_CONFIG_PATH) up 

metric-start:
	go run ./cmd/goMetricService.go -config=$(GO_CONFIG)

# Combine build and run targets
.PHONY: docker-build docker-run metric-start 
