DOCKER_CONFIG_PATH := docker/docker-compose.yaml
APP_PATH := ./cmd
BUILD_CONTEXT := .
GO_CONFIG := ./config/local.yaml
MIGRATIONS_PATH := ./migrations
CLICKHOUSE_CONNECTION = clickhouse://localhost:9000?database=testing

.PHONY: dc
dc: 
	sudo -E docker compose -f $(DOCKER_CONFIG_PATH) up 

.PHONY: dg
dg: 
	sudo docker build -t gms-container .

.PHONY: build
build:
	go build -o $(APP_PATH)/goMetricService $(APP_PATH)/goMetricService.go 

.PHONY: run
run: 
	$(APP_PATH)/goMetricService -config=$(GO_CONFIG)

.PHONY: rgui
rgui:
	sudo docker run -p 8001:8001 redislabs/redisinsight:latest	

#create migrations with "migrate" tool
.PHONY: m-c 
m-c:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq clickhouse_init

.PHONY: m-r 
m-r: 
	migrate -path $(MIGRATIONS_PATH) -database '$(CLICKHOUSE_CONNECTION)' up

.PHONY: m-d 
m-d:
	migrate -path $(MIGRATIONS_PATH) -database '$(CLICKHOUSE_CONNECTION)' down

# Combine build and run targets
