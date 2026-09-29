#!/usr/bin/env bash
# Installs repo Git hooks into .git/hooks (symlink). Run once after clone:
#   ./.githooks/install.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOOKS_DIR="${ROOT}/.git/hooks"

if [[ ! -d "${ROOT}/.git" ]]; then
  echo "error: not a git repository: ${ROOT}" >&2
  exit 1
fi

mkdir -p "${HOOKS_DIR}"
chmod +x "${ROOT}/.githooks/pre-commit" "${ROOT}/.githooks/pre-push" \
  "${ROOT}/backend/scripts/pre-commit.sh" "${ROOT}/.githooks/install.sh"

ln -sfn "../../.githooks/pre-commit" "${HOOKS_DIR}/pre-commit"
ln -sfn "../../.githooks/pre-push" "${HOOKS_DIR}/pre-push"

echo "Installed pre-commit hook -> ${HOOKS_DIR}/pre-commit"
echo "  (checks Go formatting and runs go vet for changed backend services)"
echo "Installed pre-push hook -> ${HOOKS_DIR}/pre-push"
echo "  (builds frontend and backend separately when related files are pushed)"
echo "  Skip either hook with: SKIP_GIT_HOOKS=1 git <commit|push> ..."
