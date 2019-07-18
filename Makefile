.PHONY: install lint

install:
	go install -v

lint:
	golangci-lint run

lint-all:
	golangci-lint run --enable-all
