#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="$(cd "${script_dir}/.." && pwd)"

for fixture_case in resolved empty failed; do
  M1_FIXTURE_CASE="${fixture_case}" npm --prefix "${root}/web/eino-workbench" exec -- playwright test
done
