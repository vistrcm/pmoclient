.PHONY: install lint

install:
	go install -v

lint:
	golangci-lint run
