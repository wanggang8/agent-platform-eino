#!/usr/bin/env bash
set -euo pipefail
root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"
mkdir -p test-results
exec > >(tee test-results/toolchain-baseline.log) 2>&1

if [[ ${CI:-false} == true ]]; then
  test -n "${CI_COMMIT_SHA:-}"
  git diff --quiet
  git diff --cached --quiet
  test -z "$(git ls-files --others --exclude-standard)"
else
  printf '%s\n' 'evidence_mode=preflight; dirty local results can never be G-TOOLCHAIN PASS'
fi

bash scripts/verify_toolchain.sh
bash scripts/validate_ci_config.sh

before=$(sha256sum package-lock.json go.mod go.sum)
npm ci
after=$(sha256sum package-lock.json go.mod go.sum)
test "$before" = "$after"

npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
npm run eino-workbench:openapi-lint

GOTOOLCHAIN=local go list -mod=readonly -m all
GOTOOLCHAIN=local go mod verify
GOTOOLCHAIN=local go test ./... -count=1
GOTOOLCHAIN=local go test -race ./internal/einoapp/execution ./internal/einoapp/store/sqlite -count=1
GOTOOLCHAIN=local go test ./internal/einoapp/execution ./internal/einoapp/store/sqlite -run 'CheckPoint|Checkpoint|SQLite|Migration|Store' -count=1
GOTOOLCHAIN=local go vet ./...
GOTOOLCHAIN=local go build ./...
GOTOOLCHAIN=local go test ./internal/einoapp/architecture -run 'TestImportBoundaryDoesNotUseLegacyProjectPackages|TestImportBoundaryPreservesLayering|TestProductLayersDoNotImportEino|TestPhaseOnePackageSkeletonExists' -count=1

npm run eino-workbench:typecheck
npm run eino-workbench:test
npm run eino-workbench:stream-test
npm run eino-workbench:build
npm run eino-workbench:browser-test -- --project=desktop
bash scripts/eino_workbench_server_smoke.sh --scenario contract
git diff --check

if [[ ${CI:-false} == true ]]; then
  git diff --quiet
  git diff --cached --quiet
  test -z "$(git ls-files --others --exclude-standard)"
fi
