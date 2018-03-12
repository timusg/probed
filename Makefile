.PHONY: all
all: build test

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint

build-deps:
	dep ensure

compile:
	mkdir -p build
	go build -race -o build/probed

build: build-deps compile fmt vet lint

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint -set_exit_status $(ALL_PACKAGES)

test: build-deps fmt vet build
	ENVIRONMENT=test go test $(ALL_PACKAGES)

test-cover-html:
	@echo "mode: count" > coverage-all.out

	$(foreach pkg, $(ALL_PACKAGES),\
	ENVIRONMENT=test go test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o build/coverage.html
