#!/usr/bin/env bash
set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

if ! command -v gofmt >/dev/null 2>&1; then
  echo "error: gofmt is required — install Go 1.26 or later" >&2
  exit 1
fi
if ! command -v go >/dev/null 2>&1; then
  echo "error: go is required — install Go 1.26 or later" >&2
  exit 1
fi

SERVICES=""
FORMAT_FAILED=0
STAGED_FILES="$(mktemp -d)"
trap 'rm -rf "$STAGED_FILES"' EXIT

while IFS= read -r -d '' file; do
  [[ "$file" == *.go ]] || continue
  service="${file#backend/}"
  service="${service%%/*}"

  case "$service" in
    auth-service|user-service|notification-service|api-gateway) ;;
    *) continue ;;
  esac

  SERVICES+=" ${service} "
  staged_file="${STAGED_FILES}/${file}"
  mkdir -p "$(dirname "$staged_file")"
  git show ":$file" > "$staged_file"

  if format_diff="$(gofmt -d "$staged_file")"; then
    gofmt_status=0
  else
    gofmt_status=$?
    if [[ "$gofmt_status" -ne 1 ]]; then
      echo "error: gofmt could not check staged file: $file" >&2
      exit 1
    fi
  fi
  if [[ -n "$format_diff" ]]; then
    printf '%s\n' "error: staged Go file is not gofmt-formatted: $file" >&2
    printf '%s\n' "$format_diff" >&2
    FORMAT_FAILED=1
  fi
done < <(git diff --cached --name-only --diff-filter=ACMR -z -- backend/)

if [[ "$FORMAT_FAILED" -ne 0 ]]; then
  echo "Run gofmt on the affected file(s), then stage the changes and commit again." >&2
  exit 1
fi

for service in auth-service user-service notification-service api-gateway; do
  if [[ "$SERVICES" == *" $service "* ]]; then
    echo "[pre-commit] Running go vet for ${service}..."
    (cd "${ROOT}/backend/${service}" && go vet ./...)
  fi
done

echo "[pre-commit] Backend formatting and lint checks passed"
