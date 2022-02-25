PWD := $(shell pwd)
ifeq '${CI}' 'true'
VERSION ?= dev
else
VERSION ?= $(shell git describe --tags)
VERSIONV ?= $(shell git describe --tags | sed 's,v,,g')
endif
TAG ?= "minio/operator:$(VERSION)"
LDFLAGS ?= "-s -w -X main.Version=$(VERSION)"
TMPFILE := $(shell mktemp)
GOPATH := $(shell go env GOPATH)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

HELM_HOME=helm/operator
HELM_CRDS=$(HELM_HOME)/crds

KUSTOMIZE_HOME=resources
KUSTOMIZE_CRDS=$(KUSTOMIZE_HOME)/base/crds/

PLUGIN_HOME=kubectl-minio

LOGSEARCHAPI=logsearchapi

all: build

getdeps:
	@echo "Checking dependencies"
	@mkdir -p ${GOPATH}/bin
	@which golangci-lint 1>/dev/null || (echo "Installing golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.43.0)

verify: getdeps govet gotest lint

operator:
	@CGO_ENABLED=0 GOOS=linux go build -trimpath --ldflags $(LDFLAGS) -o minio-operator
	@docker build -t $(TAG) .

build: regen-crd verify plugin logsearchapi operator

install: all

lint:
	@echo "Running $@ check"
	@GO111MODULE=on ${GOPATH}/bin/golangci-lint run --timeout=5m --config ./.golangci.yml

govet:
	@go vet ./...

gotest:
	@go test -race ./...

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@find . -name '*.zip' | xargs rm -fv
	@rm -rf dist/

regen-crd:
	@go install -v github.com/minio/controller-tools/cmd/controller-gen@v0.4.7
	@echo "WARNING: installing our fork github.com/minio/controller-tools/cmd/controller-gen@v0.4.7"
	@echo "Any other controller-gen will cause the generated CRD to lose the volumeClaimTemplate metadata to be lost"
	@${GOPATH}/bin/controller-gen crd:maxDescLen=0,generateEmbeddedObjectMeta=true paths="./..." output:crd:artifacts:config=$(KUSTOMIZE_CRDS)
	@kustomize build resources/patch-crd > $(TMPFILE)
	@mv -f $(TMPFILE) resources/base/crds/minio.min.io_tenants.yaml
	@cp -f resources/base/crds/minio.min.io_tenants.yaml $(HELM_CRDS)/minio.min.io_tenants.yaml

regen-crd-docs:
	@which crd-ref-docs 1>/dev/null || (echo "Installing crd-ref-docs" && GO111MODULE=on go install -v github.com/elastic/crd-ref-docs@latest)
	@${GOPATH}/bin/crd-ref-docs --source-path=./pkg/apis/minio.min.io/v2 --config=docs/templates/config.yaml --renderer=asciidoctor --output-path=docs/crd.adoc --templates-dir=docs/templates/asciidoctor/

plugin: regen-crd
	@echo "Building 'kubectl-minio' binary"
	@(cd $(PLUGIN_HOME); \
		go vet ./... && \
		go test -race ./... && \
		GO111MODULE=on ${GOPATH}/bin/golangci-lint cache clean && \
		GO111MODULE=on ${GOPATH}/bin/golangci-lint run --timeout=5m --config ../.golangci.yml)

.PHONY: logsearchapi
logsearchapi:
	@echo "Building 'logsearchapi' binary"
	@(cd $(LOGSEARCHAPI); \
		go vet ./... && \
		go test -race ./... && \
		GO111MODULE=on ${GOPATH}/bin/golangci-lint cache clean && \
		GO111MODULE=on ${GOPATH}/bin/golangci-lint run --timeout=5m --config ../.golangci.yml && \
		CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w" -trimpath -o ../logsearchapi-bin )

getconsoleuiyaml:
	@echo "Getting the latest Console UI"
	@kustomize build github.com/minio/console/k8s/operator-console/base > resources/base/console-ui.yaml
	@echo "Done"

generate-code:
	@./k8s/update-codegen.sh

generate-openshift-manifests:
	@./olm.sh

release: generate-openshift-manifests
	@./helm-reindex.sh
