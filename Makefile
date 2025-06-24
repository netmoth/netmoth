SHELL:= /bin/bash
ROOT_PATH:=$(abspath $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
GO_PATH:=$(shell go env GOPATH)
CPU_ARCH:=$(shell go env GOARCH)
OS_NAME:=$(shell go env GOHOSTOS)

# Оптимизированные флаги компиляции для производительности
GO_FLAGS := -ldflags="-s -w" -gcflags="-l=4" -trimpath
GO_OPTIMIZE_FLAGS := -ldflags="-s -w -extldflags=-Wl,-z,relro,-z,now" -gcflags="-l=4 -B -N" -trimpath
GO_EBPF_FLAGS := -ldflags="-s -w" -gcflags="-l=4" -trimpath -tags=ebpf

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
	openssl req -new -text -passout pass:abcd -subj /CN=Netmoth -out docker/postgres/ca/server.req -keyout docker/postgres/ca/privkey.pem
	openssl rsa -in docker/postgres/ca/privkey.pem -passin pass:abcd -out docker/postgres/ca/server.key
	openssl req -x509 -in docker/postgres/ca/server.req -text -key docker/postgres/ca/server.key -out docker/postgres/ca/server.crt
	chmod 600 docker/postgres/ca/server.key
	sudo chown 70 docker/postgres/ca/server.key
#############################################################################


#############################################################################
.PHONY: build
build: ## build
	$(eval NAME=$(filter-out $@,$(MAKECMDGOALS)))
	@if [ ${NAME} ];then\
		if [ -d ${ROOT_PATH}/cmd/${NAME}/ ];then\
			go build $(GO_FLAGS) -o bin/${NAME} cmd/${NAME}/main.go;\
		else \
			echo "error";\
		fi \
	else \
		for entry in ${ROOT_PATH}/cmd/*/;do\
			go build $(GO_FLAGS) -o bin/$$(basename $${entry}) cmd/$$(basename $${entry})/main.go;\
		done;\
	fi
#############################################################################

#############################################################################
.PHONY: build-optimized
build-optimized: ## build with maximum optimizations
	$(eval NAME=$(filter-out $@,$(MAKECMDGOALS)))
	@if [ ${NAME} ];then\
		if [ -d ${ROOT_PATH}/cmd/${NAME}/ ];then\
			go build $(GO_OPTIMIZE_FLAGS) -o bin/${NAME} cmd/${NAME}/main.go;\
		else \
			echo "error";\
		fi \
	else \
		for entry in ${ROOT_PATH}/cmd/*/;do\
			go build $(GO_OPTIMIZE_FLAGS) -o bin/$$(basename $${entry}) cmd/$$(basename $${entry})/main.go;\
		done;\
	fi
#############################################################################

#############################################################################
.PHONY: build-ebpf
build-ebpf: ## build with eBPF support
	$(eval NAME=$(filter-out $@,$(MAKECMDGOALS)))
	@if [ ${NAME} ];then\
		if [ -d ${ROOT_PATH}/cmd/${NAME}/ ];then\
			go build $(GO_EBPF_FLAGS) -o bin/${NAME} cmd/${NAME}/main.go;\
		else \
			echo "error";\
		fi \
	else \
		for entry in ${ROOT_PATH}/cmd/*/;do\
			go build $(GO_EBPF_FLAGS) -o bin/$$(basename $${entry}) cmd/$$(basename $${entry})/main.go;\
		done;\
	fi
#############################################################################

#############################################################################
.PHONY: build-ebpf-optimized
build-ebpf-optimized: ## build with eBPF support and maximum optimizations
	$(eval NAME=$(filter-out $@,$(MAKECMDGOALS)))
	@if [ ${NAME} ];then\
		if [ -d ${ROOT_PATH}/cmd/${NAME}/ ];then\
			go build $(GO_OPTIMIZE_FLAGS) -tags=ebpf -o bin/${NAME} cmd/${NAME}/main.go;\
		else \
			echo "error";\
		fi \
	else \
		for entry in ${ROOT_PATH}/cmd/*/;do\
			go build $(GO_OPTIMIZE_FLAGS) -tags=ebpf -o bin/$$(basename $${entry}) cmd/$$(basename $${entry})/main.go;\
		done;\
	fi
#############################################################################

#############################################################################
.PHONY: build-race
build-race: ## build with race detector
	$(eval NAME=$(filter-out $@,$(MAKECMDGOALS)))
	@if [ ${NAME} ];then\
		if [ -d ${ROOT_PATH}/cmd/${NAME}/ ];then\
			go build -race -o bin/${NAME} cmd/${NAME}/main.go;\
		else \
			echo "error";\
		fi \
	else \
		for entry in ${ROOT_PATH}/cmd/*/;do\
			go build -race -o bin/$$(basename $${entry}) cmd/$$(basename $${entry})/main.go;\
		done;\
	fi
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
#############################################################################


#############################################################################
%: ## A parameter
	@true
#############################################################################