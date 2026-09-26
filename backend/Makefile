BACKEND_SERVICES := auth-service user-service api-gateway

.PHONY: security-check
# Scan every backend service for known Go vulnerabilities.
security-check:
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "error: govulncheck is required — install it with:" >&2; \
		echo "  go install golang.org/x/vuln/cmd/govulncheck@latest" >&2; \
		exit 1; \
	fi
	@if ! command -v jq >/dev/null 2>&1; then \
		echo "error: jq is required to format vulnerability scan results" >&2; \
		exit 1; \
	fi
	@tmp_dir=$$(mktemp -d) || exit 1; \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	for service in $(BACKEND_SERVICES); do \
		echo "Scanning $$service..."; \
		if (cd "$$service" && govulncheck -format=json ./...) > "$$tmp_dir/$$service.json"; then \
			echo 0 > "$$tmp_dir/$$service.exit"; \
		else \
			echo $$? > "$$tmp_dir/$$service.exit"; \
		fi; \
	done; \
	bash scripts/vulnerability_table.sh "$$tmp_dir"
