VERSION := 1.0.0

GROUP := github.com/console-dns/server
SVC_NAME := console-dns
GOPATH := $(shell go env GOPATH)

DEBUG ?= false

ifdef PROD
LD_FLAGS :=  \
	-X $(GROUP)/cmd.DefConfigPath=/etc/$(SVC_NAME)/server.yaml \
	-X $(GROUP)/pkg/content/settings.ConfigLog=/var/log/$(SVC_NAME)/logs.csv \
	-X $(GROUP)/pkg/content/settings.ConfigSession=/var/lib/$(SVC_NAME)/sessions.yaml \
	-X $(GROUP)/pkg/content/settings.ConfigToken=/var/lib/$(SVC_NAME)/clients.yaml \
	-X $(GROUP)/pkg/content/settings.ConfigZone=/var/lib/$(SVC_NAME)/zones.yaml
endif


console-dns: $(shell find ./ -type f -name '*.go')
	@go build -ldflags "$(LD_FLAGS)" -o console-dns ./

dev: conf/server.yaml
	@DEBUG=true go run -tags dev ./ server -c conf/server.yaml

run: conf/server.yaml
	@go run ./ server -c conf/server.yaml

conf/server.yaml:
	@go run ./ config generate -c conf/server.yaml

fmt:
	@(test -f "$(GOPATH)/bin/gofumpt" || go install golang.org/x/tools/cmd/goimports@latest) && \
	"$(GOPATH)/bin/gofumpt" -l -w .