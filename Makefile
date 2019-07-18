.PHONY: install
install:
	go install -v

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-all
lint-all:
	golangci-lint run --enable-all

.PHONY: fmt
fmt:
	go fmt ./...
