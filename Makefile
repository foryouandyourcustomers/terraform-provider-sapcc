TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
HOST=registry.terraform.io
NAMESPACE=foryouandyourcustomers
NAME=sapcc
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
WIREMOCK=rodolpheche/wiremock:2.30.1

default: build
all: lint testacc build docs install


pull-mock:
	@docker pull ${WIREMOCK}

run-mock:
	@docker run --rm -p 8080:8080 --name wiremock -v ${PWD}/mocks:/home/wiremock ${WIREMOCK} --verbose --global-response-templating --local-response-templating

start-mock:
	@docker run --rm -d -p 8080:8080 --name wiremock -v ${PWD}/mocks:/home/wiremock ${WIREMOCK} --verbose --global-response-templating --local-response-templating
	@echo "Mock SAP Commerce Api Server available at http://localhost:8080"

stop-mock:
	@docker kill wiremock

restart-mock: stop-mock start-mock

clean:
	@echo "Cleaning up binaries"
	go clean -testcache
	rm -fr ./bin
	rm -fr ./dist
	rm -fr ./docs

build:
	@go mod tidy
	@go mod vendor
	@mkdir -p bin
	go build -o bin/${BINARY}

release:
	@goreleaser release --rm-dist --snapshot

install: build
	@mkdir -p ~/.terraform.d/plugins/${HOST}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	@rm -f ~/.terraform.d/plugins/${HOST}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}
	cp ./bin/${BINARY} ~/.terraform.d/plugins/${HOST}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}


testacc: install
	@go clean -testcache && go test  -v ./internal/provider/ -timeout=10m -count=1

fmt:
	@goimports -w $(GOFMT_FILES)

lint:
	@golangci-lint run

docs:
	@go generate

.PHONY: clean start-sapcc-mock stop-sapcc-mock restart-mock website fmt docs lint pull-mock
