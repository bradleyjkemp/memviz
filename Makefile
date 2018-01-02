.PHONY: install
install: get_dependencies install_linters

.PHONY: get_dependencies
get_dependencies:
	go get github.com/golang/dep/cmd/dep
	$(GOPATH)/bin/dep ensure

.PHONY: install_linters
install_linters:
	go get github.com/mattn/goveralls
	go get github.com/alecthomas/gometalinter
	$(GOPATH)/bin/gometalinter --install

.PHONY: test
test:
	go test ./...
	$(GOPATH)/bin/gometalinter --vendor --fast ./...

.PHONY: test-ci
test-ci: lint-ci
	$(GOPATH)/bin/goveralls -v -service=travis-ci

.PHONY: lint-ci
lint-ci:
	$(GOPATH)/bin/gometalinter --vendor --disable=gas --disable=gocyclo ./...
