TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
HOST=registry.terraform.io
NAMESPACE=fyayc
NAME=sapcc
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_amd64

default: build
all: lint test build docs install


run-mock:
	@docker run --rm -p 8080:8080 --name wiremock -v ${PWD}/mocks:/home/wiremock rodolpheche/wiremock --verbose --global-response-templating --local-response-templating

start-mock:
	@docker run --rm -d -p 8080:8080 --name wiremock -v ${PWD}/mocks:/home/wiremock rodolpheche/wiremock --verbose --global-response-templating --local-response-templating
	@echo "Mock SAP Commerce Api Server available at http://localhost:8080"

stop-mock:
	@docker kill wiremock

restart-mock: stop-mock start-mock

clean:
	@echo "Cleaning up binaries"
	rm -fr ./bin

build:
	@go mod tidy
	@go mod vendor
	@mkdir -p bin
	go build -o bin/${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	@mkdir -p ~/.terraform.d/plugins/${HOST}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	@rm -f ~/.terraform.d/plugins/${HOST}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}
	@cp ./bin/${BINARY} ~/.terraform.d/plugins/${HOST}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	@go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1  go test -v -cover ./internal/provider/

fmt:
	@goimports -w $(GOFMT_FILES)

lint:
	@golangci-lint run

docs:
	@go generate

.PHONY: clean start-sapcc-mock stop-sapcc-mock restart-mock website fmt docs lint
