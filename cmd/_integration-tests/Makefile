NO_COLOR=\e[0m
OK_COLOR=\e[38;5;118m
UNDER=\n________________________________________________________________________________\n
END_COLOR_LINE=$(UNDER)$(NO_COLOR)

TRANSPORT_TEST_MSG=\n$(OK_COLOR)Starting transport end to end test:$(END_COLOR_LINE)

MIDDLEWARES_TEST_MSG=\n$(OK_COLOR)Starting middlewares end to end test:$(END_COLOR_LINE)

CLI_TEST_MSG=\n$(OK_COLOR)Start server and cliclient generate, build, and run test:$(END_COLOR_LINE)

SHA := $(shell git rev-parse --short=10 HEAD)

MAKEFILE_PATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
VERSION_DATE := $(shell $(MAKEFILE_PATH)/../../commit_date.sh)

export PATH := $(CURDIR):$(PATH)

all: test

test: clean test-transport test-middlewares test-cli

truss:
	go install -ldflags '-X "main.Version=$(SHA)" -X "main.VersionDate=$(VERSION_DATE)"' github.com/TuneLab/truss/cmd/truss

protoc-gen-truss-protocast:
	go build -o protoc-gen-truss-protocast github.com/TuneLab/truss/cmd/protoc-gen-truss-protocast

test-transport: protoc-gen-truss-protocast truss
	@which truss
	@printf '$(TRANSPORT_TEST_MSG)'
	$(MAKE) -C transport
	rm -f ./truss ./protoc-gen-truss-protocast

test-middlewares: protoc-gen-truss-protocast truss
	@which truss
	@printf '$(MIDDLEWARES_TEST_MSG)'
	$(MAKE) -C middlewares
	rm -f ./truss ./protoc-gen-truss-protocast

test-cli: protoc-gen-truss-protocast truss
	@which truss
	@printf '$(CLI_TEST_MSG)'
	go test -v ./cli
	rm -f ./truss ./protoc-gen-truss-protocast

clean:
	rm -f ./truss ./protoc-gen-truss-protocast
	go test ./cli -clean
	$(MAKE) -C transport clean
	$(MAKE) -C middlewares clean

