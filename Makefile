# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: tidy add-copyright format lint cover build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/tianrandailove/peitho
VERSION_PACKAGE=github.com/tianrandailove/peitho/pkg/version

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/image.mk
include scripts/make-rules/copyright.mk
include scripts/make-rules/release.mk
include scripts/make-rules/dependencies.mk
include scripts/make-rules/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/build.multiarch
               Example: make build BINS="peitho"
  IMAGES       Backend images to make. Default is all of cmd starting with peitho.
               This option is available when using: make image/image.multiarch/push/push.multiarch
               Example: make image.multiarch IMAGES="peitho"
  PLATFORMS    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
               This option is available when using: make build.multiarch/image.multiarch/push.multiarch
               Example: make image.multiarch IMAGES="peitho" PLATFORMS="linux_amd64 linux_arm64"
  VERSION      The version information compiled into binaries.
               The default is obtained from gsemver or git.
  V            Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Build docker images for multiple platforms and push images to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## deploy: Deploy updated components to development env.
.PHONY: deploy
deploy:
	@$(MAKE) deploy.run

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
.PHONY: cover 
cover:
	@$(MAKE) go.test.cover

.PHONY: release.build
release.build:
	@$(MAKE) push.multiarch

## release: Release peitho
.PHONY: release
release:
	@$(MAKE) release.run

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt

## verify-copyright: Verify the boilerplate headers for all files.
.PHONY: verify-copyright
verify-copyright:
	@$(MAKE) copyright.verify

## add-copyright: Ensures source code files have copyright license headers.
.PHONY: add-copyright
add-copyright:
	@$(MAKE) copyright.add

## dependencies: Install necessary dependencies.
.PHONY: dependencies
dependencies:
	@$(MAKE) dependencies.run

## tools: install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
