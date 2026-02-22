#!/usr/bin/env bash
set -euo pipefail

# bench-ci.sh - Run benchmarks in a format compatible with benchstat.
# Outputs only go test -bench lines (no headers/comments) so benchstat can parse them.
# Usage: bash scripts/bench-ci.sh <output-file> [count]

OUT_FILE="${1:?usage: bench-ci.sh <output-file> [count]}"
COUNT="${2:-5}"

echo "running benchmarks (count=${COUNT}) -> ${OUT_FILE}"

go test ./cmd \
  -run '^$' -bench '^BenchmarkRunVersionOnlyFlag$' \
  -benchmem -count="${COUNT}" \
  > "${OUT_FILE}" 2>&1

go test ./internal/asc \
  -run '^$' -bench '^BenchmarkClientNewRequest$' \
  -benchmem -count="${COUNT}" \
  >> "${OUT_FILE}" 2>&1

go test ./internal/asc \
  -run '^$' -bench '^BenchmarkPaginate(AllAggregation|EachStreaming)$' \
  -benchmem -count="${COUNT}" \
  >> "${OUT_FILE}" 2>&1

go test ./internal/workflow \
  -run '^$' -bench '^Benchmark(BuildEnvSlice|ResolveShellCached)$' \
  -benchmem -count="${COUNT}" \
  >> "${OUT_FILE}" 2>&1

go test ./internal/screenshots \
  -run '^$' -bench '^BenchmarkEnsurePinnedKoubouVersionCached$' \
  -benchmem -count="${COUNT}" \
  >> "${OUT_FILE}" 2>&1

go test ./internal/cli/migrate \
  -run '^$' -bench '^BenchmarkInferScreenshotDisplayTypeFromDimensions$' \
  -benchmem -count="${COUNT}" \
  >> "${OUT_FILE}" 2>&1

echo "benchmarks complete: ${OUT_FILE}"
