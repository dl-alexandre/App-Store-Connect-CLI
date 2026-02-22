#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${1:-${ROOT_DIR}/.perf}"
mkdir -p "${OUT_DIR}"

STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
OUT_FILE="${OUT_DIR}/bench-${STAMP}.txt"

echo "writing benchmark snapshot to ${OUT_FILE}"
{
  echo "# asc performance snapshot"
  echo "# timestamp: ${STAMP}"
  echo "# go: $(go version)"
  echo "# host: $(uname -a)"
  echo

  echo "## startup"
  go test ./cmd -run '^$' -bench '^BenchmarkRunVersionOnlyFlag$' -benchmem -count=7
  echo

  echo "## api-path"
  go test ./internal/asc -run '^$' -bench '^BenchmarkClientNewRequest$' -benchmem -count=7
  go test ./internal/asc -run '^$' -bench '^BenchmarkPaginate(AllAggregation|EachStreaming)$' -benchmem -count=5
  echo

  echo "## local-tooling"
  go test ./internal/workflow -run '^$' -bench '^Benchmark(BuildEnvSlice|ResolveShellCached)$' -benchmem -count=7
  go test ./internal/screenshots -run '^$' -bench '^BenchmarkEnsurePinnedKoubouVersionCached$' -benchmem -count=7
  go test ./internal/cli/migrate -run '^$' -bench '^BenchmarkInferScreenshotDisplayTypeFromDimensions$' -benchmem -count=7
} | tee "${OUT_FILE}"

echo "snapshot complete: ${OUT_FILE}"
