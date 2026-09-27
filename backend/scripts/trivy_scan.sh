#!/usr/bin/env bash
set -uo pipefail

services=(auth-service user-service api-gateway)
severity="UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL"
scan_id="trivy-$$-${RANDOM}"
failed=0
images=()

usage() {
  echo "usage: $0 [all|fs|image]" >&2
}

mode="${1:-all}"
case "$mode" in
  all|fs|image) ;;
  *)
    usage
    exit 2
    ;;
esac

if ! command -v trivy >/dev/null 2>&1; then
  echo "error: Trivy is required — install it from https://trivy.dev/latest/getting-started/installation/" >&2
  exit 1
fi

if [[ "$mode" == "all" || "$mode" == "image" ]]; then
  if ! command -v docker >/dev/null 2>&1; then
    echo "error: Docker is required to build images for Trivy image scans" >&2
    exit 1
  fi
  if ! docker info >/dev/null 2>&1; then
    echo "error: the Docker daemon is unavailable; start Docker before scanning service images" >&2
    exit 1
  fi
fi

backend_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cleanup() {
  for image in "${images[@]}"; do
    if ! docker image rm "$image" >/dev/null; then
      echo "warning: could not remove temporary scan image $image" >&2
    fi
  done
}

if [[ "$mode" == "all" || "$mode" == "fs" ]]; then
  for service in "${services[@]}"; do
    echo "Scanning ${service} Go dependencies with Trivy..."
    if ! (
      cd "${backend_dir}/${service}" &&
        trivy fs --scanners vuln --severity "$severity" --exit-code 1 --format table .
    ); then
      failed=1
    fi
  done
fi

if [[ "$mode" == "all" || "$mode" == "image" ]]; then
  trap cleanup EXIT
  for service in "${services[@]}"; do
    image="backend-trivy/${service}:${scan_id}"

    echo "Building ${service} image for Trivy..."
    if ! docker build --pull --file "${backend_dir}/${service}/Dockerfile" --tag "$image" "${backend_dir}/${service}"; then
      echo "error: could not build ${service} image for Trivy scanning" >&2
      failed=1
      continue
    fi
    images+=("$image")

    echo "Scanning ${service} image with Trivy..."
    if ! trivy image --scanners vuln --severity "$severity" --exit-code 1 --format table "$image"; then
      failed=1
    fi
  done
fi

exit "$failed"
