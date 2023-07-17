OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
BINARY_PATH = $$(cat dist/artifacts.json | jq '.[] | select(.goarch == "$(ARCH)" and .goos == "$(OS)") | .path' -r)
GOBIN_PATH = $$PWD/.bin
ENV_VARS = GOBIN="$(GOBIN_PATH)" PATH="$(GOBIN_PATH):$$PATH"
.PHONY: release
release: tools
	@$(ENV_VARS) goreleaser release $(RELEASE_ARGS)
.PHONY: build
build: tools
	@$(ENV_VARS) goreleaser build --snapshot --rm-dist
	@ln -sf $(BINARY_PATH) ./gh-terraport
.PHONY: install
install: build
	@gh extension remove gh-terraport || true
	@gh extension install --force .
.PHONY: tools
tools:
	$(ENV_VARS) go install $$(go list -f '{{join .Imports " "}}' tools.go)