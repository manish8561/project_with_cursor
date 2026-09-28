BACKEND_SERVICES := auth-service user-service notification-service api-gateway

.PHONY: test security-check trivy-security-check trivy-fs-security-check trivy-image-security-check

# Run unit tests for every backend service.
test:
	@set -e; \
	for service in $(BACKEND_SERVICES); do \
		echo "Testing $$service..."; \
		(cd "$$service" && go test -v ./...); \
	done

# Scan every backend service for known Go and image vulnerabilities.
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
	bash scripts/vulnerability_table.sh "$$tmp_dir"; govulncheck_status=$$?; \
	bash scripts/trivy_scan.sh all; trivy_status=$$?; \
	test "$$govulncheck_status" -eq 0 -a "$$trivy_status" -eq 0

# Run Trivy without the govulncheck scan.
trivy-security-check:
	bash scripts/trivy_scan.sh all

trivy-fs-security-check:
	bash scripts/trivy_scan.sh fs

trivy-image-security-check:
	bash scripts/trivy_scan.sh image
