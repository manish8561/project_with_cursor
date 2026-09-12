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
chmod +x "${ROOT}/.githooks/pre-push" "${ROOT}/.githooks/install.sh"

# Remove old pre-commit hook from earlier setup if present.
if [[ -L "${HOOKS_DIR}/pre-commit" ]] || [[ -f "${HOOKS_DIR}/pre-commit" ]]; then
  rm -f "${HOOKS_DIR}/pre-commit"
fi

ln -sfn "../../.githooks/pre-push" "${HOOKS_DIR}/pre-push"

echo "Installed pre-push hook -> ${HOOKS_DIR}/pre-push"
echo "  (builds frontend and backend separately when related files are pushed)"
echo "  Skip with: SKIP_GIT_HOOKS=1 git push ..."
