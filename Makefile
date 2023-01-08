# MAKEFILE FOR DOAN

VERSION=$(shell cat .release-please-manifest.json | jq -r '."."')
BINARY_NAME=doan
GOOS=linux
GOARCH=amd64
TARGET=cmd/agent/main.go
BINARY_OUTPUT="bin/$(GOOS)/$(GOARCH)/$(BINARY_NAME)"
REPO_URL=https://github.com/mjmorales/doan


.PHONY: build-binary
build-binary:
	@echo "Building $(BINARY_OUTPUT) with version $(VERSION) "
	@go build -o $(BINARY_OUTPUT) -ldflags "-X github.com/mjmorales/doan/pkg/agent.version=$(VERSION)" $(TARGET)

.PHONY: build-deb
build-deb: build-binary
	@echo "Building deb package"
	rm -rf build && mkdir -p build
	fpm \
		-s dir -t deb \
		-C ./bin/$(GOOS)/$(GOARCH) \
		-p release/$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).deb \
		--name $(BINARY_NAME) \
		--version $(VERSION) \
		--architecture $(GOARCH) \
		--description "Daemon for handling scheduled ansible runs." \
		--url "https://github.com/mjmorales/doan" \
		--maintainer "Manuel Morales <morales.jmanuel16@gmail.com>" \
		doan=/usr/local/bin/doan

.PHONY: release
release:
	@echo "Creating debian changelog for version $(VERSION)"
	release-please release-pr \
		--repo-url=$(REPO_URL) \
		--token=${GH_TOKEN} \
		--release-type=go \
		--package-name=$(BINARY_NAME)
