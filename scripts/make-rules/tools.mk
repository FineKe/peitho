# ==============================================================================
# Makefile helper functions for tools
#

TOOLS ?=$(BLOCKER_TOOLS) $(CRITICAL_TOOLS) $(TRIVIAL_TOOLS)

.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
	@golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi

.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) get github.com/jstemmer/go-junit-report
	@$(GO) install github.com/jstemmer/go-junit-report

.PHONY: install.gsemver
install.gsemver:
	@$(GO) get github.com/arnaud-deprez/gsemver
	@$(GO) install github.com/arnaud-deprez/gsemver

.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) get github.com/git-chglog/git-chglog/cmd/git-chglog
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog

.PHONY: install.github-release
install.github-release:
	@$(GO) get github.com/github-release/github-release
	@$(GO) install github.com/github-release/github-release

.PHONY: install.golines
install.golines:
	@$(GO) get github.com/segmentio/golines
	@$(GO) install github.com/segmentio/golines

.PHONY: install.go-mod-outdated
install.go-mod-outdated:
	@$(GO) get github.com/psampaz/go-mod-outdated
	@$(GO) install github.com/psampaz/go-mod-outdated

.PHONY: install.mockgen
install.mockgen:
	@$(GO) get github.com/golang/mock/mockgen
	@$(GO) install github.com/golang/mock/mockgen

.PHONY: install.gotests
install.gotests:
	@$(GO) get github.com/cweill/gotests/...
	@$(GO) install github.com/cweill/gotests/...

.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) get github.com/golang/protobuf/protoc-gen-go
	@$(GO) install github.com/golang/protobuf/protoc-gen-go

.PHONY: install.addlicense
install.addlicense:
	@$(GO) get github.com/marmotedu/addlicense
	@$(GO) install github.com/marmotedu/addlicense

.PHONY: install.goimports
install.goimports:
	@$(GO) get golang.org/x/tools/cmd/goimports
	@$(GO) install golang.org/x/tools/cmd/goimports

.PHONY: install.depth
install.depth:
	@$(GO) get github.com/KyleBanks/depth/cmd/depth
	@$(GO) install github.com/KyleBanks/depth/cmd/depth

.PHONY: install.go-callvis
install.go-callvis:
	@$(GO) get github.com/ofabry/go-callvis
	@$(GO) install github.com/ofabry/go-callvis

.PHONY: install.gothanks
install.gothanks:
	@$(GO) get github.com/psampaz/gothanks
	@$(GO) install github.com/psampaz/gothanks

.PHONY: install.richgo
install.richgo:
	@$(GO) get github.com/kyoh86/richgo
	@$(GO) install github.com/kyoh86/richgo
