package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const validFixture = `name: Toolchain Gate

on:
  push:
  workflow_dispatch:

permissions: {}

jobs:
  toolchain-build:
    runs-on: ubuntu-24.04
    permissions:
      contents: read
      packages: write
    outputs:
      toolchain_image: ${{ steps.build.outputs.TOOLCHAIN_IMAGE }}
    steps:
      - name: Checkout
        uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0
        with:
          fetch-depth: 0
      - name: Setup Docker
        uses: docker/setup-docker-action@6d7cfa65f60a9dda7b46e5513fa982536f3c9877
        with:
          version: v29.4.0
      - name: Configure memory-only GHCR credentials
        run: |
          install -d -m 0700 "$RUNNER_TEMP/docker-config" "$RUNNER_TEMP/docker-credential-bin"
          printf '%s\n' '{"credHelpers":{"ghcr.io":"github-token"}}' > "$RUNNER_TEMP/docker-config/config.json"
          install -m 0700 scripts/docker-credential-github-token "$RUNNER_TEMP/docker-credential-bin/docker-credential-github-token"
          printf 'DOCKER_CONFIG=%s\n' "$RUNNER_TEMP/docker-config" >> "$GITHUB_ENV"
          printf '%s\n' "$RUNNER_TEMP/docker-credential-bin" >> "$GITHUB_PATH"
      - name: Setup Buildx
        uses: docker/setup-buildx-action@bb05f3f5519dd87d3ba754cc423b652a5edd6d2c
        with:
          version: v0.35.0
          driver-opts: image=moby/buildkit:v0.31.1@sha256:6b59b7df63a8cb9902736f9ddf7fcff8261613d3e7449b8ea8b7537fc399c03a
      - name: Build and push canonical image
        id: build
        env:
          GHCR_ACTOR: ${{ github.actor }}
          GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          repository=${GITHUB_REPOSITORY,,}
          image_ref="ghcr.io/$repository/toolchain:$GITHUB_SHA"
          bash scripts/build_toolchain_image.sh --push "$image_ref" --env-file test-results/toolchain.env
          cat test-results/toolchain.env >> "$GITHUB_OUTPUT"

  toolchain-verify:
    needs: toolchain-build
    runs-on: ubuntu-24.04
    permissions:
      contents: read
      packages: read
    env:
      TOOLCHAIN_IMAGE: ${{ needs.toolchain-build.outputs.toolchain_image }}
    steps:
      - name: Checkout
        uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0
        with:
          fetch-depth: 0
      - name: Setup Docker
        uses: docker/setup-docker-action@6d7cfa65f60a9dda7b46e5513fa982536f3c9877
        with:
          version: v29.4.0
      - name: Configure memory-only GHCR credentials
        run: |
          install -d -m 0700 "$RUNNER_TEMP/docker-config" "$RUNNER_TEMP/docker-credential-bin"
          printf '%s\n' '{"credHelpers":{"ghcr.io":"github-token"}}' > "$RUNNER_TEMP/docker-config/config.json"
          install -m 0700 scripts/docker-credential-github-token "$RUNNER_TEMP/docker-credential-bin/docker-credential-github-token"
          printf 'DOCKER_CONFIG=%s\n' "$RUNNER_TEMP/docker-config" >> "$GITHUB_ENV"
          printf '%s\n' "$RUNNER_TEMP/docker-credential-bin" >> "$GITHUB_PATH"
      - name: Verify digest image
        env:
          GHCR_ACTOR: ${{ github.actor }}
          GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          [[ "$TOOLCHAIN_IMAGE" =~ ^ghcr\.io/[a-z0-9._/-]+/toolchain:[0-9a-f]{40}@sha256:[0-9a-f]{64}$ ]]
          docker pull "$TOOLCHAIN_IMAGE"
          docker run --rm --platform linux/amd64 -e CI=true -e GITHUB_SHA="$GITHUB_SHA" -v "$GITHUB_WORKSPACE:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh
      - name: Upload toolchain evidence
        if: ${{ always() }}
        uses: actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a
        with:
          name: toolchain-evidence-${{ github.sha }}
          path: |
            test-results/toolchain-baseline.log
            test-results/eino-workbench-playwright-report/
          if-no-files-found: error
          retention-days: 30
`

func completeTestLock() lockValues {
	return lockValues{
		GitHubRunner:             "ubuntu-24.04",
		ActionsCheckoutSHA:       "9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0",
		ActionsUploadArtifactSHA: "043fb46d1a93c77aae656e7c1c64a875d1fc6a0a",
		DockerSetupDockerSHA:     "6d7cfa65f60a9dda7b46e5513fa982536f3c9877",
		DockerSetupBuildxSHA:     "bb05f3f5519dd87d3ba754cc423b652a5edd6d2c",
		DockerEngineVersion:      "29.4.0",
		DockerBuildxVersion:      "0.35.0",
		BuildKitImage:            "moby/buildkit:v0.31.1",
		BuildKitDigest:           "sha256:6b59b7df63a8cb9902736f9ddf7fcff8261613d3e7449b8ea8b7537fc399c03a",
	}
}

func TestValidateAcceptsCanonicalGitHubWorkflow(t *testing.T) {
	if err := validateConfig([]byte(validFixture), completeTestLock()); err != nil {
		t.Fatalf("expected canonical GitHub workflow to pass, got %v", err)
	}
	config := parseTestConfig(t, validFixture)
	references, err := canonicalScriptReferences(config)
	if err != nil {
		t.Fatalf("expected canonical script references, got %v", err)
	}
	want := []string{"scripts/docker-credential-github-token", "scripts/build_toolchain_image.sh", "scripts/run_toolchain_baseline.sh"}
	if !reflect.DeepEqual(references, want) {
		t.Fatalf("expected only canonical scripts %v, got %v", want, references)
	}
}

func TestValidateRejectsTriggerDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "missing push", mutate: func(config map[string]any) { delete(config["on"].(map[string]any), "push") }},
		{name: "missing workflow dispatch", mutate: func(config map[string]any) { delete(config["on"].(map[string]any), "workflow_dispatch") }},
		{name: "unexpected pull request", mutate: func(config map[string]any) { config["on"].(map[string]any)["pull_request"] = nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, "workflow triggers must be exactly push and workflow_dispatch")
		})
	}
}

func TestValidateRejectsTopLevelDrift(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		summary string
	}{
		{name: "wrong name", mutate: func(config map[string]any) { config["name"] = "Other" }, summary: "workflow name must be Toolchain Gate"},
		{name: "non-empty permissions", mutate: func(config map[string]any) { config["permissions"] = map[string]any{"contents": "read"} }, summary: "top-level permissions must be empty"},
		{name: "extra top-level key", mutate: func(config map[string]any) { config["concurrency"] = "toolchain" }, summary: "workflow top-level keys must be exact"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, tt.summary)
		})
	}
}

func TestValidateRejectsJobSetDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "missing build", mutate: func(config map[string]any) { delete(config["jobs"].(map[string]any), buildJobName) }},
		{name: "missing verify", mutate: func(config map[string]any) { delete(config["jobs"].(map[string]any), verifyJobName) }},
		{name: "extra job", mutate: func(config map[string]any) { config["jobs"].(map[string]any)["extra"] = map[string]any{} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, "jobs must be exactly toolchain-build and toolchain-verify")
		})
	}
}

func TestValidateRejectsRunnerAndPermissionDrift(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		summary string
	}{
		{name: "build runner", mutate: func(config map[string]any) { testJob(config, buildJobName)["runs-on"] = "ubuntu-latest" }, summary: "toolchain-build runner must match lock"},
		{name: "verify runner", mutate: func(config map[string]any) { testJob(config, verifyJobName)["runs-on"] = "ubuntu-latest" }, summary: "toolchain-verify runner must match lock"},
		{name: "build packages missing", mutate: func(config map[string]any) {
			delete(testJob(config, buildJobName)["permissions"].(map[string]any), "packages")
		}, summary: "toolchain-build permissions must be exact"},
		{name: "build packages read", mutate: func(config map[string]any) {
			testJob(config, buildJobName)["permissions"].(map[string]any)["packages"] = "read"
		}, summary: "toolchain-build permissions must be exact"},
		{name: "verify packages missing", mutate: func(config map[string]any) {
			delete(testJob(config, verifyJobName)["permissions"].(map[string]any), "packages")
		}, summary: "toolchain-verify permissions must be exact"},
		{name: "verify packages write", mutate: func(config map[string]any) {
			testJob(config, verifyJobName)["permissions"].(map[string]any)["packages"] = "write"
		}, summary: "toolchain-verify permissions must be exact"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, tt.summary)
		})
	}
}

func TestValidateRejectsActionPinAndVersionDrift(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		summary string
	}{
		{name: "checkout tag", mutate: func(config map[string]any) { testStep(config, buildJobName, 0)["uses"] = "actions/checkout@v7" }, summary: "Checkout action must use the lock SHA"},
		{name: "checkout branch", mutate: func(config map[string]any) { testStep(config, buildJobName, 0)["uses"] = "actions/checkout@main" }, summary: "Checkout action must use the lock SHA"},
		{name: "checkout different SHA", mutate: func(config map[string]any) {
			testStep(config, verifyJobName, 0)["uses"] = "actions/checkout@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}, summary: "Checkout action must use the lock SHA"},
		{name: "setup docker SHA", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 1)["uses"] = "docker/setup-docker-action@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}, summary: "Setup Docker action must match lock"},
		{name: "setup buildx SHA", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 3)["uses"] = "docker/setup-buildx-action@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}, summary: "Setup Buildx action must match lock"},
		{name: "engine version", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 1)["with"].(map[string]any)["version"] = "v29.3.0"
		}, summary: "Setup Docker action must match lock"},
		{name: "buildx version", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 3)["with"].(map[string]any)["version"] = "v0.34.0"
		}, summary: "Setup Buildx action must match lock"},
		{name: "buildkit image", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 3)["with"].(map[string]any)["driver-opts"] = "image=moby/buildkit:v0.30.0@" + completeTestLock().BuildKitDigest
		}, summary: "Setup Buildx action must match lock"},
		{name: "buildkit digest", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 3)["with"].(map[string]any)["driver-opts"] = "image=" + completeTestLock().BuildKitImage + "@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}, summary: "Setup Buildx action must match lock"},
		{name: "artifact SHA", mutate: func(config map[string]any) {
			testStep(config, verifyJobName, 4)["uses"] = "actions/upload-artifact@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}, summary: "artifact step must match lock"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, tt.summary)
		})
	}
}

func TestValidateRejectsBuildBoundaryDrift(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		summary string
	}{
		{name: "output source", mutate: func(config map[string]any) {
			testJob(config, buildJobName)["outputs"].(map[string]any)["toolchain_image"] = "${{ steps.other.outputs.TOOLCHAIN_IMAGE }}"
		}, summary: "toolchain-build output must use steps.build.outputs.TOOLCHAIN_IMAGE"},
		{name: "builder target", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 4)["run"] = strings.Replace(testStep(config, buildJobName, 4)["run"].(string), `ghcr.io/$repository/toolchain:$GITHUB_SHA`, `ghcr.io/$repository/toolchain:latest`, 1)
		}, summary: "Build and push canonical image step must match canonical command"},
		{name: "builder id", mutate: func(config map[string]any) { testStep(config, buildJobName, 4)["id"] = "other" }, summary: "Build and push canonical image step must match canonical command"},
		{name: "extra step", mutate: func(config map[string]any) {
			job := testJob(config, buildJobName)
			job["steps"] = append(job["steps"].([]any), map[string]any{"run": "echo extra"})
		}, summary: "toolchain-build steps must be exact"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, tt.summary)
		})
	}
}

func TestValidateRejectsCredentialBoundaryDrift(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		summary string
	}{
		{name: "docker login regression", mutate: func(config map[string]any) {
			testStep(config, buildJobName, 2)["run"] = "docker login ghcr.io"
		}, summary: "credential helper configuration must match canonical command"},
		{name: "config contains token", mutate: func(config map[string]any) {
			step := testStep(config, buildJobName, 2)
			step["run"] = strings.Replace(step["run"].(string), `{"credHelpers":{"ghcr.io":"github-token"}}`, `{"token":"$GHCR_TOKEN"}`, 1)
		}, summary: "credential helper configuration must match canonical command"},
		{name: "helper name drift", mutate: func(config map[string]any) {
			step := testStep(config, verifyJobName, 2)
			step["run"] = strings.ReplaceAll(step["run"].(string), "github-token", "other-helper")
		}, summary: "credential helper configuration must match canonical command"},
		{name: "DOCKER_CONFIG drift", mutate: func(config map[string]any) {
			step := testStep(config, buildJobName, 2)
			step["run"] = strings.Replace(step["run"].(string), "DOCKER_CONFIG=%s", "DOCKER_CONFIG_FILE=%s", 1)
		}, summary: "credential helper configuration must match canonical command"},
		{name: "GITHUB_ENV drift", mutate: func(config map[string]any) {
			step := testStep(config, buildJobName, 2)
			step["run"] = strings.Replace(step["run"].(string), "$GITHUB_ENV", "$GITHUB_PATH", 1)
		}, summary: "credential helper configuration must match canonical command"},
		{name: "GITHUB_PATH drift", mutate: func(config map[string]any) {
			step := testStep(config, verifyJobName, 2)
			step["run"] = strings.Replace(step["run"].(string), "$GITHUB_PATH", "$GITHUB_ENV", 1)
		}, summary: "credential helper configuration must match canonical command"},
		{name: "build token at job scope", mutate: func(config map[string]any) {
			testJob(config, buildJobName)["env"] = map[string]any{"GHCR_TOKEN": "${{ secrets.GITHUB_TOKEN }}"}
		}, summary: "toolchain-build keys must be exact"},
		{name: "verify token at job scope", mutate: func(config map[string]any) {
			testJob(config, verifyJobName)["env"].(map[string]any)["GHCR_TOKEN"] = "${{ secrets.GITHUB_TOKEN }}"
		}, summary: "toolchain-verify image must come from toolchain-build output"},
		{name: "build missing actor", mutate: func(config map[string]any) {
			delete(testStep(config, buildJobName, 4)["env"].(map[string]any), "GHCR_ACTOR")
		}, summary: "Build and push canonical image authentication must be step-local"},
		{name: "build missing token", mutate: func(config map[string]any) {
			delete(testStep(config, buildJobName, 4)["env"].(map[string]any), "GHCR_TOKEN")
		}, summary: "Build and push canonical image authentication must be step-local"},
		{name: "verify missing actor", mutate: func(config map[string]any) {
			delete(testStep(config, verifyJobName, 3)["env"].(map[string]any), "GHCR_ACTOR")
		}, summary: "Verify digest image authentication must be step-local"},
		{name: "verify missing token", mutate: func(config map[string]any) {
			delete(testStep(config, verifyJobName, 3)["env"].(map[string]any), "GHCR_TOKEN")
		}, summary: "Verify digest image authentication must be step-local"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, tt.summary)
		})
	}
}

func TestValidateRejectsVerifyBoundaryDrift(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		summary string
	}{
		{name: "missing need", mutate: func(config map[string]any) { delete(testJob(config, verifyJobName), "needs") }, summary: "toolchain-verify must need toolchain-build"},
		{name: "wrong need", mutate: func(config map[string]any) { testJob(config, verifyJobName)["needs"] = "other" }, summary: "toolchain-verify must need toolchain-build"},
		{name: "tag-only guard", mutate: func(config map[string]any) {
			step := testStep(config, verifyJobName, 3)
			step["run"] = strings.Replace(step["run"].(string), `@sha256:[0-9a-f]{64}`, "", 1)
		}, summary: "Verify digest image step must match canonical command"},
		{name: "pull command", mutate: func(config map[string]any) {
			step := testStep(config, verifyJobName, 3)
			step["run"] = strings.Replace(step["run"].(string), `docker pull "$TOOLCHAIN_IMAGE"`, `docker pull ghcr.io/other/image:latest`, 1)
		}, summary: "Verify digest image step must match canonical command"},
		{name: "baseline command", mutate: func(config map[string]any) {
			step := testStep(config, verifyJobName, 3)
			step["run"] = strings.Replace(step["run"].(string), "bash scripts/run_toolchain_baseline.sh", "bash scripts/other.sh", 1)
		}, summary: "Verify digest image step must match canonical command"},
		{name: "verify image source", mutate: func(config map[string]any) {
			testJob(config, verifyJobName)["env"].(map[string]any)["TOOLCHAIN_IMAGE"] = "${{ needs.other.outputs.toolchain_image }}"
		}, summary: "toolchain-verify image must come from toolchain-build output"},
		{name: "extra step", mutate: func(config map[string]any) {
			job := testJob(config, verifyJobName)
			job["steps"] = append(job["steps"].([]any), map[string]any{"run": "echo extra"})
		}, summary: "toolchain-verify steps must be exact"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, tt.summary)
		})
	}
}

func TestValidateRejectsArtifactDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "path", mutate: func(config map[string]any) {
			testStep(config, verifyJobName, 4)["with"].(map[string]any)["path"] = "test-results/other.log"
		}},
		{name: "retention", mutate: func(config map[string]any) {
			testStep(config, verifyJobName, 4)["with"].(map[string]any)["retention-days"] = 29
		}},
		{name: "missing always", mutate: func(config map[string]any) { delete(testStep(config, verifyJobName, 4), "if") }},
		{name: "conditional success", mutate: func(config map[string]any) { testStep(config, verifyJobName, 4)["if"] = "${{ success() }}" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			assertConfigError(t, config, "artifact step must match lock")
		})
	}
}

func TestValidateRejectsGitLabSyntax(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "GitLab variable", text: strings.Replace(validFixture, "$GITHUB_SHA", "$CI_COMMIT_SHA", 1)},
		{name: "GitLab config", text: validFixture + "# .gitlab-ci.yml\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig([]byte(tt.text), completeTestLock())
			if err == nil || !strings.Contains(err.Error(), "GitLab syntax is forbidden") {
				t.Fatalf("expected stable GitLab syntax error, got %v", err)
			}
		})
	}
}

func TestValidateRejectsMalformedYAML(t *testing.T) {
	err := validateConfig([]byte("jobs: ["), completeTestLock())
	if err == nil || !strings.Contains(err.Error(), "invalid YAML") {
		t.Fatalf("expected YAML parse error, got %v", err)
	}
}

func TestLoadLockAcceptsCanonicalPins(t *testing.T) {
	path := writeTestLock(t, canonicalTestLock())
	got, err := loadLock(path)
	if err != nil {
		t.Fatalf("expected lock to load, got %v", err)
	}
	if !reflect.DeepEqual(got, completeTestLock()) {
		t.Fatalf("loaded lock mismatch: %#v", got)
	}
}

func TestLoadLockAcceptsPinnedBuildKitRepositoryWithExactSemver(t *testing.T) {
	content := strings.Replace(canonicalTestLock(), completeTestLock().BuildKitImage, "ghcr.io/example/buildkit:v1.2.3", 1)
	got, err := loadLock(writeTestLock(t, content))
	if err != nil {
		t.Fatalf("expected lock-owned BuildKit repository to load, got %v", err)
	}
	if got.BuildKitImage != "ghcr.io/example/buildkit:v1.2.3" {
		t.Fatalf("expected alternate pinned BuildKit image, got %q", got.BuildKitImage)
	}
}

func TestLoadLockRejectsMissingDuplicateAndMalformedEntries(t *testing.T) {
	tests := []struct {
		name    string
		content string
		summary string
	}{
		{name: "missing key", content: strings.Replace(canonicalTestLock(), "GITHUB_RUNNER=ubuntu-24.04\n", "", 1), summary: "missing toolchain lock key GITHUB_RUNNER"},
		{name: "missing runtime key", content: strings.Replace(canonicalTestLock(), "PLAYWRIGHT_VERSION=1.61.1\n", "", 1), summary: "missing toolchain lock key PLAYWRIGHT_VERSION"},
		{name: "unknown key", content: canonicalTestLock() + "UNKNOWN_PIN=unexpected\n", summary: "unknown toolchain lock key UNKNOWN_PIN"},
		{name: "duplicate key", content: canonicalTestLock() + "GITHUB_RUNNER=ubuntu-24.04\n", summary: "duplicate toolchain lock key GITHUB_RUNNER"},
		{name: "malformed key", content: strings.Replace(canonicalTestLock(), "GITHUB_RUNNER=", "github_runner=", 1), summary: "invalid toolchain lock entry"},
		{name: "whitespace", content: strings.Replace(canonicalTestLock(), "ubuntu-24.04", "ubuntu 24.04", 1), summary: "invalid toolchain lock entry"},
		{name: "shell syntax", content: strings.Replace(canonicalTestLock(), "ubuntu-24.04", "$(uname)", 1), summary: "invalid toolchain lock entry"},
		{name: "short action SHA", content: strings.Replace(canonicalTestLock(), completeTestLock().ActionsCheckoutSHA, "abc123", 1), summary: "invalid ACTIONS_CHECKOUT_SHA pin"},
		{name: "uppercase action SHA", content: strings.Replace(canonicalTestLock(), completeTestLock().ActionsCheckoutSHA, strings.ToUpper(completeTestLock().ActionsCheckoutSHA), 1), summary: "invalid ACTIONS_CHECKOUT_SHA pin"},
		{name: "invalid semver", content: strings.Replace(canonicalTestLock(), "DOCKER_ENGINE_VERSION=29.4.0", "DOCKER_ENGINE_VERSION=v29.4", 1), summary: "invalid DOCKER_ENGINE_VERSION pin"},
		{name: "invalid buildkit image", content: strings.Replace(canonicalTestLock(), completeTestLock().BuildKitImage, "example/buildkit:latest", 1), summary: "invalid BUILDKIT_IMAGE pin"},
		{name: "invalid digest", content: strings.Replace(canonicalTestLock(), completeTestLock().BuildKitDigest, "sha256:abc", 1), summary: "invalid BUILDKIT_DIGEST pin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadLock(writeTestLock(t, tt.content))
			if err == nil || !strings.Contains(err.Error(), tt.summary) {
				t.Fatalf("expected %q, got %v", tt.summary, err)
			}
		})
	}
}

func TestValidateScriptReferencesRejectsEscapes(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"docker-credential-github-token", "build_toolchain_image.sh", "run_toolchain_baseline.sh"} {
		if err := os.WriteFile(filepath.Join(root, "scripts", name), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := validateScriptReferences(root, parseTestConfig(t, validFixture)); err != nil {
		t.Fatalf("expected canonical references to remain inside root, got %v", err)
	}

	outside := filepath.Join(t.TempDir(), "outside.sh")
	if err := os.WriteFile(outside, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	// credential helper 也必须参与真实路径检查，不能通过 symlink 逃逸仓库边界。
	helper := filepath.Join(root, "scripts", "docker-credential-github-token")
	if err := os.Remove(helper); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, helper); err != nil {
		t.Fatal(err)
	}
	err := validateScriptReferences(root, parseTestConfig(t, validFixture))
	if err == nil || !strings.Contains(err.Error(), "script reference escapes repository") {
		t.Fatalf("expected credential helper symlink escape error, got %v", err)
	}
	if err := os.Remove(helper); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(helper, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		target string
	}{
		{name: "absolute", target: outside},
		{name: "parent traversal", target: "../outside.sh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			step := testStep(config, buildJobName, 4)
			step["run"] = strings.Replace(step["run"].(string), "scripts/build_toolchain_image.sh", tt.target, 1)
			err := validateScriptReferences(root, config)
			if err == nil || !strings.Contains(err.Error(), "script reference escapes repository") {
				t.Fatalf("expected script escape error, got %v", err)
			}
		})
	}

	if err := os.Symlink(outside, filepath.Join(root, "scripts", "escape.sh")); err != nil {
		t.Fatal(err)
	}
	config := parseTestConfig(t, validFixture)
	step := testStep(config, buildJobName, 4)
	step["run"] = strings.Replace(step["run"].(string), "scripts/build_toolchain_image.sh", "scripts/escape.sh", 1)
	err = validateScriptReferences(root, config)
	if err == nil || !strings.Contains(err.Error(), "script reference escapes repository") {
		t.Fatalf("expected symlink escape error, got %v", err)
	}
}

func canonicalTestLock() string {
	lock := completeTestLock()
	return strings.Join([]string{
		"PLATFORM=linux/amd64",
		"PLAYWRIGHT_IMAGE=mcr.microsoft.com/playwright:v1.61.1-noble",
		"PLAYWRIGHT_AMD64_DIGEST=sha256:cf0daee9b994042e011bc29f20cdff1a9f682a039b43fcd738f7d8a9d3bcd9d6",
		"PLAYWRIGHT_VERSION=1.61.1",
		"CHROMIUM_REVISION=1228",
		"CHROMIUM_VERSION=149.0.7827.55",
		"BASE_OS=ubuntu-24.04-noble",
		"FONT_POLICY=playwright-v1.61.1-noble-bundled",
		"UBUNTU_SNAPSHOT=20260708T000000Z",
		"APT_BUILD_PACKAGES=build-essential",
		"GO_LINUX_AMD64_SHA256=5c2c3b16caefa1d968a94c1daca04a7ca301a496d9b086e17ad77bb81393f053",
		"NODE_LINUX_X64_SHA256=783130984963db7ba9cbd01089eaf2c2efb055c7c1693c943174b967b3050cb8",
		"GITHUB_RUNNER=" + lock.GitHubRunner,
		"ACTIONS_CHECKOUT_SHA=" + lock.ActionsCheckoutSHA,
		"ACTIONS_UPLOAD_ARTIFACT_SHA=" + lock.ActionsUploadArtifactSHA,
		"DOCKER_SETUP_DOCKER_SHA=" + lock.DockerSetupDockerSHA,
		"DOCKER_SETUP_BUILDX_SHA=" + lock.DockerSetupBuildxSHA,
		"DOCKER_ENGINE_VERSION=" + lock.DockerEngineVersion,
		"DOCKER_BUILDX_VERSION=" + lock.DockerBuildxVersion,
		"BUILDKIT_IMAGE=" + lock.BuildKitImage,
		"BUILDKIT_DIGEST=" + lock.BuildKitDigest,
	}, "\n") + "\n"
}

func writeTestLock(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "toolchain.lock")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func assertConfigError(t *testing.T, config map[string]any, summary string) {
	t.Helper()
	err := validateConfig(marshalTestConfig(t, config), completeTestLock())
	if err == nil || !strings.Contains(err.Error(), summary) {
		t.Fatalf("expected %q, got %v", summary, err)
	}
}

func testJob(config map[string]any, name string) map[string]any {
	return config["jobs"].(map[string]any)[name].(map[string]any)
}

func testStep(config map[string]any, jobName string, index int) map[string]any {
	return testJob(config, jobName)["steps"].([]any)[index].(map[string]any)
}

func parseTestConfig(t *testing.T, fixture string) map[string]any {
	t.Helper()
	var config map[string]any
	if err := yaml.Unmarshal([]byte(fixture), &config); err != nil {
		t.Fatal(err)
	}
	return config
}

func marshalTestConfig(t *testing.T, config map[string]any) []byte {
	t.Helper()
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
