.PHONY: install
install: get_dependencies install_linters

.PHONY: get_dependencies
get_dependencies:
	go get github.com/golang/dep/cmd/dep
	$(GOPATH)/bin/dep ensure

.PHONY: install_linters
install_linters:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.12.5

.PHONY: test
test: lint
	go test ./...

.PHONY: test-ci
test-ci: lint
	$(GOPATH)/bin/goveralls -v -service=travis-ci

.PHONY: lint
lint:
	$(GOPATH)/bin/golangci-lint run
