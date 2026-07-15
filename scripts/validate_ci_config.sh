#!/usr/bin/env bash
set -euo pipefail
root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"
exec env GOTOOLCHAIN=local go run ./scripts/validate_ci_config \
  -config .github/workflows/toolchain.yml -lock build/toolchain/toolchain.lock -root "$root"
