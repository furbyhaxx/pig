# SPDX-License-Identifier: MIT

# Fork gates (SPEC-0001 D10). Upstream's Pi-comparison, porting-ledger and publication gates stay defined in the other makefiles but are not part of this set.

FORK_BASE ?= @{upstream}

##@ Fork

fork-check: build lint-changed fork-vet-windows test-changed go-fix-clean divergence-guard port-lint perf-lint any-trap source-hygiene docs-drift changelog-check ## Fork gate before a commit or a ready report (SPEC-0001 D10)
	@echo "make fork-check: all fork gates green."

# lint-changed already runs govet natively; Windows is a shipped target, so vet the directly changed packages for it too.
fork-vet-windows: ## go vet the directly changed packages for windows/amd64 (FORK_BASE, default @{upstream})
	@packages=$$(python3 automation/ci/changed-packages.py --base '$(FORK_BASE)' | awk '$$1 == "direct" { print $$2 }') || exit 1; \
	if [ -z "$$packages" ]; then echo "fork-vet-windows: no changed Go packages"; exit 0; fi; \
	GOOS=windows GOARCH=amd64 go vet $$packages

.PHONY: fork-check fork-vet-windows
