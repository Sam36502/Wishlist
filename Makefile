EXE_LINUX = "run_server_linux"
EXE_WIN = "run_server_win.exe"
DOCKER_IMAGE = "wishlist"
CONTAINER_NAME ="wishlist"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Cleans up generated target files
	@echo "### Cleaning up ###"
	rm -f data/${EXE_LINUX} data/${EXE_WIN}

build: ## Builds the executable for linux
	@echo "### Building Linux Executable ###"
	@GOOS="linux" CGO_ENABLED=0 go build -o data/${EXE_LINUX} ./src/

build-win: ## Builds the executable for windows
	@echo "### Building Windows Executable ###"
	@GOOS="windows" go build -o data/${EXE_WIN} ./src/

image: build ## Builds the docker image
	@echo "### Building Docker Image ###"
	@docker build -t ${DOCKER_IMAGE} .

up: down image ## Starts the container
	@echo "### Starting Container ###"
	@docker run -d --name ${CONTAINER_NAME} -v "/etc/letsencrypt:/certs:ro" -p 80:80 ${DOCKER_IMAGE}

down: ## Stops the container
	@echo "### Stopping Container ###"
	@-docker stop ${CONTAINER_NAME}
	@-docker rm ${CONTAINER_NAME}
