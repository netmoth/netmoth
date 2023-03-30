SHELL:= /bin/bash
ROOT_PATH:=$(abspath $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
GO_PATH:=$(shell go env GOPATH)
CPU_ARCH:=$(shell go env GOARCH)
OS_NAME:=$(shell go env GOHOSTOS)

NAME:=netmoth

.DEFAULT_GOAL:=help

#############################################################################
.PHONY: help
help: ## help
	@grep --no-filename -E '^[a-zA-Z_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
#############################################################################


#############################################################################
.PHONY: init
init: ## init
	mkdir docker/postgres docker/postgres/ca
	openssl req -new -text -passout pass:abcd -subj /CN=Werbot -out docker/postgres/ca/server.req -keyout docker/postgres/ca/privkey.pem
	openssl rsa -in docker/postgres/ca/privkey.pem -passin pass:abcd -out docker/postgres/ca/server.key
	openssl req -x509 -in docker/postgres/ca/server.req -text -key docker/postgres/ca/server.key -out docker/postgres/ca/server.crt
	chmod 600 docker/postgres/ca/server.key
	sudo chown 70 docker/postgres/ca/server.key
#############################################################################


#############################################################################
.PHONY: run
run: build ## run
	@sudo ./bin/${NAME}
#############################################################################


#############################################################################
.PHONY: build
build: ## build
	@go build -o bin/${NAME} cmd/agent/main.go
#	@upx -1 -k bin/${NAME}
#	@rm -rf bin/${NAME}.*
#############################################################################


#############################################################################
.PHONY: lint
lint: ## lint
	@$(GO_PATH)/bin/golangci-lint run
#############################################################################


#############################################################################
.PHONY: clean
clean: ## clean
	@rm -rf $(ROOT_PATH)/bin/*
	@rm -rf $(ROOT_PATH)/*.log
	@rm -rf $(ROOT_PATH)/cmd/*.log
#############################################################################