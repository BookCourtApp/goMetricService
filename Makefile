DOCKER_CONFIG_PATH := Docker/docker-compose.yaml
BUILD_CONTEXT := .

# Build the Docker container
docker-build:
	sudo docker build -t clickhouse -f $(DOCKERFILE_PATH) $(BUILD_CONTEXT)

docker-run:
	sudo docker compose -f $(DOCKER_CONFIG_PATH) up

# Combine build and run targets
.PHONY: docker-build docker-run 
