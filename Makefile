install: lint
	go install -v

lint:
	gometalinter --vendor ./...
