PLUGIN_NAME := cm-push

HAS_PIP := $(shell command -v pip3;)
HAS_VENV := $(shell command -v virtualenv;)

.PHONY: build
build: build_linux build_linux_arm64 build_mac build_mac_arm64 build_windows build_windows_arm64

build_windows: export GOARCH=amd64
build_windows: export GO111MODULE=on
build_windows:
	@GOOS=windows go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/windows/amd64/helm-cm-push cmd/helm-cm-push/main.go  # windows

build_windows_arm64: export GOARCH=arm64
build_windows_arm64: export GO111MODULE=on
build_windows_arm64:
	@GOOS=windows go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/windows/arm64/helm-cm-push cmd/helm-cm-push/main.go  # windows arm64

link_windows:
	@cp bin/windows/amd64/helm-cm-push ./bin/helm-cm-push

link_windows_arm64:
	@cp bin/windows/arm64/helm-cm-push ./bin/helm-cm-push

build_linux: export GOARCH=amd64
build_linux: export CGO_ENABLED=0
build_linux: export GO111MODULE=on
build_linux:
	@GOOS=linux go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/linux/amd64/helm-cm-push cmd/helm-cm-push/main.go  # linux

build_linux_arm64: export GOARCH=arm64
build_linux_arm64: export CGO_ENABLED=0
build_linux_arm64: export GO111MODULE=on
build_linux_arm64:
	@GOOS=linux go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/linux/arm64/helm-cm-push cmd/helm-cm-push/main.go  # linux arm64

link_linux:
	@cp bin/linux/amd64/helm-cm-push ./bin/helm-cm-push

link_linux_arm64:
	@cp bin/linux/arm64/helm-cm-push ./bin/helm-cm-push

build_mac: export GOARCH=amd64
build_mac: export CGO_ENABLED=0
build_mac: export GO111MODULE=on
build_mac:
	@GOOS=darwin go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/darwin/amd64/helm-cm-push cmd/helm-cm-push/main.go # mac osx intel
	@cp bin/darwin/amd64/helm-cm-push ./bin/helm-cm-push # For use w make install

build_mac_arm64: export GOARCH=arm64
build_mac_arm64: export CGO_ENABLED=0
build_mac_arm64: export GO111MODULE=on
build_mac_arm64:
	@GOOS=darwin go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/darwin/arm64/helm-cm-push cmd/helm-cm-push/main.go # mac osx apple silicon
	@cp bin/darwin/arm64/helm-cm-push ./bin/helm-cm-push # For use w make install

link_mac:
	@cp bin/darwin/amd64/helm-cm-push ./bin/helm-cm-push

link_mac_arm64:
	@cp bin/darwin/arm64/helm-cm-push ./bin/helm-cm-push

.PHONY: clean
clean:
	@git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: test
test: setup-test-environment
	@./scripts/test.sh

.PHONY: covhtml
covhtml:
	@go tool cover -html=.cover/cover.out

.PHONY: tree
tree:
	@tree -I vendor

.PHONY: release
release:
	@scripts/release.sh $(VERSION)

.PHONY: install
install:
	HELM_PUSH_PLUGIN_NO_INSTALL_HOOK=1 helm plugin install $(shell pwd)

.PHONY: remove
remove:
	helm plugin remove $(PLUGIN_NAME)

.PHONY: setup-test-environment
setup-test-environment:
ifndef HAS_PIP
	@apt-get update && apt-get install -y python-pip
endif
ifndef HAS_VENV
	@pip install virtualenv
endif
	@./scripts/setup_test_environment.sh

.PHONY: acceptance
acceptance: setup-test-environment
	@./scripts/acceptance.sh
