install: lint
	go install -v

lint:
	gometalinter --vendor ./...


prereq:
	dep ensure -v && \
	go get -u github.com/alecthomas/gometalinter && \
	gometalinter --install
