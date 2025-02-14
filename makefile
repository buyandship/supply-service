HOST=bravo-lu
ENGINE=docker

.PHONY: build_all build docker_build

build_all: build docker_build

build:
	@echo "Building go binary..."
	@mkdir -p output/
	@GOOS=linux GOARCH=amd64 go build ./
	@echo "Build successfully"

docker_build:
	@echo "Building docker image..."
	@${ENGINE} build -f dockerfile -t ghcr.io/${HOST}/supply-svr:latest .
	@echo "Pushing docker image..."
	@${ENGINE} push ghcr.io/${HOST}/supply-svr:latest

