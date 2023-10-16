DOCKER_CONFIG_PATH := docker/docker-compose.yaml
PROJECT_NAME := test_uml_project
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

#UML handling, required: goplantuml, plantuml for ubuntu, gio(system linux tool)
.PHONY: u-create
u-create:
	mkdir -p UML
	goplantuml -aggregate-private-members -show-aggregations -recursive -title="$(PROJECT_NAME)" ./ > ./UML/ClassDiagram.puml	

.PHONY: u-generate
u-generate:
	plantuml -tpng -output images ./UML/*.puml

.PHONY: u-open
u-open:
	gio open UML/images/*.png


.PHONY: uml
uml: u-create u-generate u-open
