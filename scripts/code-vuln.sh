#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
export GOTOOLCHAIN=local
export GOWORK=off
module=${1:-server}
case "$module" in
  "server") scanner=v1.7.0 ;;
  *) echo "Unknown module: $module" >&2; exit 2 ;;
esac
cd "$module"
expected=$(awk '$1 == "toolchain" { sub(/^go/, "", $2); toolchain=$2 } $1 == "go" { minimum=$2 } END { print toolchain ? toolchain : minimum }' go.mod)
actual=$(go env GOVERSION)
if [[ "$actual" != "go$expected" ]]; then
  echo "Install Go $expected and put its bin directory first on PATH (found $actual). Automatic switching is disabled." >&2
  exit 2
fi
printf 'Scanning %s with %s, govulncheck %s (GOTOOLCHAIN=local)\n' "$module" "$actual" "$scanner"
# Versioned go run does not add the scanner to application dependencies.
go run "golang.org/x/vuln/cmd/govulncheck@$scanner" ./...
