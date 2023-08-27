DOCKER_CONFIG_PATH := docker/docker-compose.yaml
APP_PATH := ./cmd
BUILD_CONTEXT := .
GO_CONFIG := ./config/local.yaml

# Build the Docker container
docker-build:
	sudo docker build -t clickhouse -f $(DOCKERFILE_PATH) $(BUILD_CONTEXT)

docker-run:
	sudo docker compose -f $(DOCKER_CONFIG_PATH) up 

build:
	go build -o $(APP_PATH)/goMetricService $(APP_PATH)/goMetricService.go 

run:
	$(APP_PATH)/goMetricService -config=$(GO_CONFIG)

redis-gui:
	sudo docker run -p 8001:8001 redislabs/redisinsight:latest	

# Combine build and run targets
.PHONY: docker-build docker-run build run redis-guid
