# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
# Constants
SERVER_BINARY_PATH=dist/server
SERVER_MAIN_PATH=cmd/server/main.go
MIGRATOR_BINARY_PATH=dist/migrator
MIGRATOR_MAIN_PATH=cmd/migrator/main.go
SETUP_BINARY_PATH=dist/setup
SETUP_MAIN_PATH=cmd/setup/main.go
MANAGER_BINARY_PATH=dist/manager
MANAGER_MAIN_PATH=cmd/manager/main.go
DOCKER_COMPOSE_YML=docker-compose.yml
DOCKER_IMAGE_NAME=FlagField-Server

ifeq ($(OS), Windows_NT)
    PLATFORM="Windows"
    SERVER_BINARY_PATH=dist/server.exe
    MIGRATOR_BINARY_PATH=dist/migrator.exe
    SETUP_BINARY_PATH=dist/setup.exe
    MANAGER_BINARY_PATH=dist/manager.exe
else
    ifeq ($(shell uname), Darwin)
        PLATFORM="macOS"
    else
        PLATFORM="Linux"
    endif
endif

.PHONY: all test build clean
all: test build
build:
	$(GOBUILD) -o $(SERVER_BINARY_PATH) -ldflags "-s -w" $(SERVER_MAIN_PATH)
	$(GOBUILD) -o $(MIGRATOR_BINARY_PATH) -ldflags "-s -w" $(MIGRATOR_MAIN_PATH)
	$(GOBUILD) -o $(SETUP_BINARY_PATH) -ldflags "-s -w" $(SETUP_MAIN_PATH)
	$(GOBUILD) -o $(MANAGER_BINARY_PATH) -ldflags "-s -w" $(MANAGER_MAIN_PATH)
tools:
	$(GOBUILD) -o $(MIGRATOR_BINARY_PATH) -ldflags "-s -w" $(MIGRATOR_MAIN_PATH)
	$(GOBUILD) -o $(SETUP_BINARY_PATH) -ldflags "-s -w" $(SETUP_MAIN_PATH)
	$(GOBUILD) -o $(MANAGER_BINARY_PATH) -ldflags "-s -w" $(MANAGER_MAIN_PATH)
test:
	$(GOTEST) ./test -v
clean:
	rm -rf $(SERVER_BINARY_PATH)
	rm -rf $(MIGRATOR_BINARY_PATH)
	rm -rf $(SETUP_BINARY_PATH)
	rm -rf $(MANAGER_BINARY_PATH)
up:
	docker-compose -p $(DOCKER_IMAGE_NAME) -f $(DOCKER_COMPOSE_YML) up --build $(EXTRA)
down:
	docker-compose -p $(DOCKER_IMAGE_NAME) -f $(DOCKER_COMPOSE_YML) down
start:
	docker-compose -p $(DOCKER_IMAGE_NAME) -f $(DOCKER_COMPOSE_YML) start
stop:
	docker-compose -p $(DOCKER_IMAGE_NAME) -f $(DOCKER_COMPOSE_YML) stop
