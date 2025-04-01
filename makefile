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

.PHONY: t2
t2:
	@echo "Restarting docker service on t2"
	@ssh -i ~/.ssh/bns_mkp_dev.pem ec2-user@mkp-ssh2.buynship.com 'bash /home/ec2-user/build_supply_svr.sh'
	@echo "Script execution completed on remote server"