.PHONY: install
install: install_linters

.PHONY: install_linters
install_linters:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.12.5

.PHONY: test
test: lint
	go test ./...

.PHONY: test-ci
test-ci: lint
	go run github.com/mattn/goveralls -v -service=travis-ci

.PHONY: lint
lint:
	$(GOPATH)/bin/golangci-lint run
