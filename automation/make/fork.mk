# SPDX-License-Identifier: MIT

# Fork gates (SPEC-0001 D10). Upstream's Pi-comparison, porting-ledger and publication gates stay defined in the other makefiles but are not part of this set.

FORK_BASE ?= @{upstream}

##@ Fork

fork-check: build lint-changed fork-test go-fix-clean divergence-guard port-lint perf-lint any-trap source-hygiene docs-drift changelog-check ## Fork gate before a push: cheap static gates plus the tests of directly edited packages (SPEC-0001 D10)
	@echo "make fork-check: all fork gates green."

# Only packages the change edits; dependents are not fanned out (that is make test-changed, run on request).
fork-test: ## go test -race the directly edited Go packages (FORK_BASE, default @{upstream})
	@packages=$$(python3 automation/ci/changed-packages.py --base '$(FORK_BASE)' | awk '$$1 == "direct" { print $$2 }') || exit 1; \
	if [ -z "$$packages" ]; then echo "fork-test: no changed Go packages"; exit 0; fi; \
	echo "fork-test: $$(printf '%s' "$$packages" | tr '\n' ' ')"; \
	# internal/codingagent alone exceeds the 10-minute default under -race.
	go test -race -short -timeout 25m $$packages

.PHONY: fork-check fork-test
