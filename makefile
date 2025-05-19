.PHONY: build_all build docker_build docker_login

build_all: build docker_build

docker_login:
	aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 211125742375.dkr.ecr.ap-southeast-1.amazonaws.com


build:
	@echo "Building go binary..."
	@mkdir -p output/
	@GOOS=linux GOARCH=amd64 go build ./
	@echo "Build successfully"

docker_build:
	@echo "Building docker image..."
	@docker build -f dockerfile -t 211125742375.dkr.ecr.ap-southeast-1.amazonaws.com/supply-service:latest .
	@echo "Pushing docker image..."
	@docker push 211125742375.dkr.ecr.ap-southeast-1.amazonaws.com/supply-service:latest

.PHONY: t2
t2:
	@echo "Restarting docker service on t2"
	@ssh -i ~/.ssh/bns_mkp_dev.pem ec2-user@mkp-ssh2.buynship.com 'bash /home/ec2-user/build_supply_svr.sh'
	@echo "Script execution completed on remote server"

.PHONY: update_thrift
update_thrift:
	@hz update -idl supply.thrift