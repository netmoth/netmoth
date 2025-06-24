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
.PHONY: build-agent
build-agent: ## build agent only
	go build $(GO_FLAGS) -o bin/agent cmd/agent/main.go
#############################################################################

#############################################################################
.PHONY: build-manager
build-manager: ## build manager (central server) only
	go build $(GO_FLAGS) -o bin/manager cmd/manager/main.go
#############################################################################

#############################################################################
.PHONY: run-agent
run-agent: ## run agent with default config
	@if [ ! -f bin/agent ]; then \
		echo "Agent not built. Run 'make build-agent' first."; \
		exit 1; \
	fi
	@if [ ! -f cmd/agent/config.yml ]; then \
		echo "Agent config not found. Please create cmd/agent/config.yml"; \
		echo "Available configs: config.yml.example, config_optimized.yml, config_ebpf.yml"; \
		exit 1; \
	fi
	./scripts/run_agent.sh
#############################################################################

#############################################################################
.PHONY: run-agent-optimized
run-agent-optimized: ## run agent with optimized config
	@if [ ! -f bin/agent ]; then \
		echo "Agent not built. Run 'make build-agent' first."; \
		exit 1; \
	fi
	cp cmd/agent/config_optimized.yml cmd/agent/config.yml
	./scripts/run_agent.sh
#############################################################################

#############################################################################
.PHONY: run-agent-ebpf
run-agent-ebpf: ## run agent with eBPF config
	@if [ ! -f bin/agent ]; then \
		echo "Agent not built. Run 'make build-agent' first."; \
		exit 1; \
	fi
	cp cmd/agent/config_ebpf.yml cmd/agent/config.yml
	./scripts/run_agent.sh
#############################################################################

#############################################################################
.PHONY: run-manager
run-manager: ## run manager (central server)
	@if [ ! -f bin/manager ]; then \
		echo "Manager not built. Run 'make build-manager' first."; \
		exit 1; \
	fi
	@if [ ! -f cmd/manager/config.yml ]; then \
		echo "Manager config not found. Please create cmd/manager/config.yml"; \
		echo "Available configs: config.yml.example, config_optimized.yml, config_ebpf.yml"; \
		exit 1; \
	fi
	./bin/manager
#############################################################################

#############################################################################
.PHONY: run-manager-optimized
run-manager-optimized: ## run manager with optimized config
	@if [ ! -f bin/manager ]; then \
		echo "Manager not built. Run 'make build-manager' first."; \
		exit 1; \
	fi
	cp cmd/manager/config_optimized.yml cmd/manager/config.yml
	./bin/manager
#############################################################################

#############################################################################
.PHONY: run-manager-ebpf
run-manager-ebpf: ## run manager with eBPF config
	@if [ ! -f bin/manager ]; then \
		echo "Manager not built. Run 'make build-manager' first."; \
		exit 1; \
	fi
	cp cmd/manager/config_ebpf.yml cmd/manager/config.yml
	./bin/manager
#############################################################################

#############################################################################
.PHONY: deploy-agent
deploy-agent: ## deploy agent to remote machine (requires SSH access)
	@echo "Usage: make deploy-agent HOST=user@hostname CONFIG=cmd/agent/config.yml"
	@if [ -z "$(HOST)" ]; then \
		echo "Error: HOST parameter is required"; \
		echo "Example: make deploy-agent HOST=user@192.168.1.100"; \
		exit 1; \
	fi
	@if [ ! -f bin/agent ]; then \
		echo "Agent not built. Run 'make build-agent' first."; \
		exit 1; \
	fi
	@if [ ! -f "$(CONFIG)" ]; then \
		echo "Config file $(CONFIG) not found."; \
		exit 1; \
	fi
	scp bin/agent $(CONFIG) scripts/run_agent.sh $(HOST):~/netmoth/
	ssh $(HOST) "cd ~/netmoth && chmod +x run_agent.sh"
	@echo "Agent deployed to $(HOST)"
	@echo "To run: ssh $(HOST) 'cd ~/netmoth && sudo ./run_agent.sh'"
#############################################################################

#############################################################################
%: ## A parameter
	@true
#############################################################################