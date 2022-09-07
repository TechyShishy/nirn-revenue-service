VERSION=$(shell git describe --tags --always --long --dirty)
WINDOWS=windows
LINUX=linux
DARWIN=darwin

.PHONY: all test clean $(WINDOWS) $(LINUX) $(DARWIN) pbgos

.SUFFIXES: .go

all: test build ## Build and run tests

test: ## Run unit tests

build: windows linux darwin ## Build binaries
	@echo version: $(VERSION)

clean: ## Remove previous build
	find -iname \*.pb.go -delete
	rm -rf bin/


help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

PCMDS=cmd/%/main.go
CMDS=cmd/*/main.go
WINDOWS_POBJS=bin/%_$(WINDOWS).exe
LINUX_POBJS=bin/%_$(LINUX)
DARWIN_POBJS=bin/%_$(DARWIN)

WINDOWS_OBJECTS=$(patsubst $(PCMDS),$(WINDOWS_POBJS),$(wildcard $(CMDS)))
LINUX_OBJECTS=$(patsubst $(PCMDS),$(LINUX_POBJS),$(wildcard $(CMDS)))
DARWIN_OBJECTS=$(patsubst $(PCMDS),$(DARWIN_POBJS),$(wildcard $(CMDS)))

$(WINDOWS_POBJS): $(PCMDS)
	env GOOS=windows GOARCH=amd64 go build -o $@ -v -ldflags="-s -w -X main.version=$(VERSION)"  $<

$(LINUX_POBJS): $(PCMDS)
	env GOOS=linux GOARCH=amd64 go build -o $@ -v -ldflags="-s -w -X main.version=$(VERSION)"  $<

$(DARWIN_POBJS): $(PCMDS)
	env GOOS=darwin GOARCH=amd64 go build -o $@ -v -ldflags="-s -w -X main.version=$(VERSION)" $<

$(WINDOWS): $(WINDOWS_OBJECTS)
$(LINUX): $(LINUX_OBJECTS)
$(DARWIN): $(DARWIN_OBJECTS)

PPROTOS=%.proto
PPBGOS=%.pb.go
PBGOS=$(patsubst $(PPROTOS),$(PPBGOS),$(wildcard api/proto/*.proto))


$(PPBGOS): $(PPROTOS)
	protoc --go_out=. --go_opt=paths=source_relative $<

$(CMDS): $(PBGOS)