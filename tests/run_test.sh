#!/usr/bin/env bash
# Local CI. Lives in <project>/tests/run_tests.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="$ROOT_DIR/tests/test_outputs"
rm -rf $OUT_DIR
mkdir -p "$OUT_DIR"

# Load the root .env
set -a
source "$ROOT_DIR/.env"
set +a

G=$'\e[0;32m'
R=$'\e[0;31m'
NC=$'\e[0m'

export GOTOOLCHAIN="${GOTOOLCHAIN:-go1.27.1}"
export PATH="$(go env GOPATH)/bin:$PATH"

# '--reset' flag to reset installations
GOBIN_DIR="$(go env GOPATH)/bin"
if [[ $* == '--re-install' ]]; then
    echo "Remove old versions"
    rm -f "$GOBIN_DIR/govulncheck" "$GOBIN_DIR/golangci-lint"
    echo "Re-installations"
    go install golang.org/x/vuln/cmd/govulncheck@latest
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
    echo "${G}DONE${NC}"
    echo ""
else
    echo "run '--re-install' to install right version of golang tools"
    echo ""
fi

run() {
  local title="$1" log="$2"
  shift 2
  echo "RUN $title"
  if "$@" > "$OUT_DIR/$log" 2>&1; then
    echo "${G}DONE${NC}"
  else
    echo "${R}FAILED${NC} (see tests/test_outputs/$log)"
  fi
  echo ""
}

run "Govulncheck" govulncheck.log bash -c "cd '$ROOT_DIR' && '$GOBIN_DIR/govulncheck' ./..."

run "Golangci" golangci.log bash -c "cd '$ROOT_DIR' && '$GOBIN_DIR/golangci-lint' run"

run "Build" build.log bash -c "cd '$ROOT_DIR' && go mod download && go build -v ./..."
run "Unit tests" unit_tests.log bash -c "cd '$ROOT_DIR' && go test -v ./..."

echo "ALL TESTS PASSED"
