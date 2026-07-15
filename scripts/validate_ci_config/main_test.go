package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const canonicalRunCommand = `docker run --rm --platform linux/amd64 -e CI=true -e CI_COMMIT_SHA="$CI_COMMIT_SHA" -v "$CI_PROJECT_DIR:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh`

// validFixture 证明两个 job 通过 extends 继承 runner 镜像、DinD 与变量后仍能通过完整 DAG 校验。
const validFixture = `stages:
  - toolchain
  - verify

.docker-amd64:
  image: docker:29.4.0-cli@sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
  services:
    - name: docker:29.4.0-dind@sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd
      alias: docker
  variables:
    DOCKER_HOST: tcp://docker:2375
    DOCKER_TLS_CERTDIR: ""
    DOCKER_BUILDKIT: "1"

toolchain-build:
  extends: .docker-amd64
  stage: toolchain
  script:
    - mkdir -p test-results
    - bash scripts/build_toolchain_image.sh --push "$CI_REGISTRY_IMAGE/toolchain:$CI_COMMIT_SHA" --env-file test-results/toolchain.env
  artifacts:
    reports:
      dotenv: test-results/toolchain.env

toolchain-verify:
  extends: .docker-amd64
  stage: verify
  needs:
    - job: toolchain-build
      artifacts: true
  script:
    - test -n "$TOOLCHAIN_IMAGE"
    - 'case "$TOOLCHAIN_IMAGE" in *@sha256:*) ;; *) exit 1 ;; esac'
    - docker pull "$TOOLCHAIN_IMAGE"
    - docker run --rm --platform linux/amd64 -e CI=true -e CI_COMMIT_SHA="$CI_COMMIT_SHA" -v "$CI_PROJECT_DIR:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh
`

const validFixtureWithoutNeed = `stages: [toolchain, verify]
.docker-amd64:
  image: docker:29.4.0-cli@sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
  services:
    - name: docker:29.4.0-dind@sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd
      alias: docker
  variables:
    DOCKER_HOST: tcp://docker:2375
    DOCKER_TLS_CERTDIR: ""
    DOCKER_BUILDKIT: "1"
toolchain-build:
  extends: .docker-amd64
  stage: toolchain
  script:
    - bash scripts/build_toolchain_image.sh --push "$CI_REGISTRY_IMAGE/toolchain:$CI_COMMIT_SHA" --env-file test-results/toolchain.env
  artifacts:
    reports:
      dotenv: test-results/toolchain.env
toolchain-verify:
  extends: .docker-amd64
  stage: verify
  script:
    - 'case "$TOOLCHAIN_IMAGE" in *@sha256:*) ;; *) exit 1 ;; esac'
    - docker pull "$TOOLCHAIN_IMAGE"
    - docker run --rm --platform linux/amd64 -e CI=true -e CI_COMMIT_SHA="$CI_COMMIT_SHA" -v "$CI_PROJECT_DIR:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh
`

func completeTestLock() lockValues {
	return lockValues{
		DockerCLIImage:   "docker:29.4.0-cli",
		DockerCLIDigest:  "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		DockerDindImage:  "docker:29.4.0-dind",
		DockerDindDigest: "sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
	}
}

func TestValidateAcceptsInheritedPinnedDAG(t *testing.T) {
	if err := validateConfig([]byte(validFixture), completeTestLock()); err != nil {
		t.Fatalf("expected inherited pinned DAG to pass, got %v", err)
	}
}

func TestValidateRejectsUnpinnedDockerImage(t *testing.T) {
	config := []byte("stages: [toolchain, verify]\ntoolchain-build:\n  image: docker:29.4.0-cli\n")
	err := validateConfig(config, lockValues{
		DockerCLIImage:  "docker:29.4.0-cli",
		DockerCLIDigest: "sha256:bb21349a52c00b206ad8b5c03fa52023c741c4cf11f269d40d40b9ebaac73d96",
	})
	if err == nil || !strings.Contains(err.Error(), "digest") {
		t.Fatalf("expected digest error, got %v", err)
	}
}

func TestValidateRejectsMissingVerifyNeed(t *testing.T) {
	config := []byte(validFixtureWithoutNeed)
	if err := validateConfig(config, completeTestLock()); err == nil || !strings.Contains(err.Error(), "needs") {
		t.Fatalf("expected missing build dependency to fail, got %v", err)
	}
}

func TestValidateRejectsMalformedYAML(t *testing.T) {
	if err := validateConfig([]byte("stages: [toolchain, verify"), completeTestLock()); err == nil || !strings.Contains(err.Error(), "YAML") {
		t.Fatalf("expected YAML parse error, got %v", err)
	}
}

func TestValidateRejectsWrongStages(t *testing.T) {
	config := strings.Replace(validFixture, "  - toolchain\n  - verify", "  - verify\n  - toolchain", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "stages") {
		t.Fatalf("expected stage order error, got %v", err)
	}
}

func TestValidateRejectsMissingVerifyStage(t *testing.T) {
	config := strings.Replace(validFixture, "  - verify\n", "", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "stages") {
		t.Fatalf("expected missing verify stage error, got %v", err)
	}
}

func TestValidateRejectsJobStageDrift(t *testing.T) {
	config := strings.Replace(validFixture, "  stage: verify", "  stage: toolchain", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "stage") {
		t.Fatalf("expected job stage error, got %v", err)
	}
}

func TestValidateRejectsUnpinnedDind(t *testing.T) {
	config := strings.Replace(validFixture, completeTestLock().DockerDindImage+"@"+completeTestLock().DockerDindDigest, completeTestLock().DockerDindImage, 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "DinD digest") {
		t.Fatalf("expected DinD digest error, got %v", err)
	}
}

func TestValidateRejectsMissingDockerVariable(t *testing.T) {
	config := strings.Replace(validFixture, "    DOCKER_BUILDKIT: \"1\"\n", "", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "DOCKER_BUILDKIT") {
		t.Fatalf("expected inherited variable error, got %v", err)
	}
}

func TestValidateRejectsMissingDotenvArtifact(t *testing.T) {
	config := strings.Replace(validFixture, "      dotenv: test-results/toolchain.env", "      dotenv: test-results/other.env", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "dotenv") {
		t.Fatalf("expected dotenv error, got %v", err)
	}
}

func TestValidateRejectsBuildWithoutCanonicalBuilder(t *testing.T) {
	config := strings.Replace(validFixture, "bash scripts/build_toolchain_image.sh", "bash scripts/other.sh", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "build_toolchain_image.sh") {
		t.Fatalf("expected canonical builder error, got %v", err)
	}
}

func TestValidateRejectsBuilderMentionWithoutInvocation(t *testing.T) {
	// 单纯把路径输出到日志不能冒充 canonical builder 的真实执行。
	config := strings.Replace(validFixture, "    - bash scripts/build_toolchain_image.sh", "    - echo bash scripts/build_toolchain_image.sh", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "build_toolchain_image.sh") {
		t.Fatalf("expected non-executed builder mention to fail, got %v", err)
	}
}

func TestValidateRejectsWrongBuilderPushTarget(t *testing.T) {
	config := strings.Replace(validFixture, `"$CI_REGISTRY_IMAGE/toolchain:$CI_COMMIT_SHA"`, `registry.example/toolchain:latest`, 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "build_toolchain_image.sh") {
		t.Fatalf("expected non-canonical push target to fail, got %v", err)
	}
}

func TestValidateRejectsNeedWithoutArtifacts(t *testing.T) {
	config := strings.Replace(validFixture, "      artifacts: true", "      artifacts: false", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "artifacts") {
		t.Fatalf("expected need artifacts error, got %v", err)
	}
}

func TestValidateRejectsVerifyWithoutDigestGuard(t *testing.T) {
	config := strings.Replace(validFixture, "    - 'case \"$TOOLCHAIN_IMAGE\" in *@sha256:*) ;; *) exit 1 ;; esac'\n", "", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "@sha256:") {
		t.Fatalf("expected TOOLCHAIN_IMAGE digest guard error, got %v", err)
	}
}

func TestValidateRejectsDigestTextWithoutGuard(t *testing.T) {
	// 出现 digest 文本不等于运行时校验，必须有针对 TOOLCHAIN_IMAGE 的 guard。
	config := strings.Replace(validFixture, "    - 'case \"$TOOLCHAIN_IMAGE\" in *@sha256:*) ;; *) exit 1 ;; esac'", "    - echo marker-@sha256:value", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "@sha256:") {
		t.Fatalf("expected non-guard digest text to fail, got %v", err)
	}
}

func TestValidateRejectsVerifyWithoutBaseline(t *testing.T) {
	config := strings.Replace(validFixture, "bash scripts/run_toolchain_baseline.sh", "bash scripts/verify_toolchain.sh", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "baseline") {
		t.Fatalf("expected baseline invocation error, got %v", err)
	}
}

func TestValidateRejectsBaselineMentionWithoutInvocation(t *testing.T) {
	config := strings.Replace(validFixture, "    - "+canonicalRunCommand, "    - echo "+canonicalRunCommand, 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "baseline") {
		t.Fatalf("expected non-executed baseline mention to fail, got %v", err)
	}
}

func TestValidateRejectsVerifyUsingAnotherImage(t *testing.T) {
	config := strings.Replace(validFixture, "docker pull \"$TOOLCHAIN_IMAGE\"", "docker pull ubuntu:latest", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "TOOLCHAIN_IMAGE") {
		t.Fatalf("expected canonical image error, got %v", err)
	}
}

func TestValidateRejectsAdditionalUnpinnedDockerRun(t *testing.T) {
	// 即使 canonical baseline 存在，也不能夹带第二个未固定镜像执行。
	config := strings.Replace(validFixture, "    - docker pull \"$TOOLCHAIN_IMAGE\"", "    - docker pull \"$TOOLCHAIN_IMAGE\"\n    - docker run --rm ubuntu:latest true", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "only TOOLCHAIN_IMAGE") {
		t.Fatalf("expected additional image run to fail, got %v", err)
	}
}

func TestValidateRejectsToolchainImageUsedAsContainerArgument(t *testing.T) {
	bypass := `docker run --rm ubuntu:latest "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh`
	config := strings.Replace(validFixture, canonicalRunCommand, bypass, 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "only TOOLCHAIN_IMAGE") {
		t.Fatalf("expected non-image TOOLCHAIN_IMAGE usage to fail, got %v", err)
	}
}

func TestLoadLockReadsDockerPins(t *testing.T) {
	path := filepath.Join(t.TempDir(), "toolchain.lock")
	data := "DOCKER_CLI_IMAGE=docker:29.4.0-cli\n" +
		"DOCKER_CLI_AMD64_DIGEST=sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n" +
		"DOCKER_DIND_IMAGE=docker:29.4.0-dind\n" +
		"DOCKER_DIND_AMD64_DIGEST=sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd\n"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := loadLock(path)
	if err != nil {
		t.Fatalf("expected lock to load, got %v", err)
	}
	if got != completeTestLock() {
		t.Fatalf("unexpected lock values: %#v", got)
	}
}

func TestLoadLockRejectsDuplicateDockerPin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "toolchain.lock")
	data := "DOCKER_CLI_IMAGE=docker:29.4.0-cli\nDOCKER_CLI_IMAGE=docker:29.4.0-cli\n"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadLock(path); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("expected duplicate lock key error, got %v", err)
	}
}

func TestValidateScriptReferencesAcceptsExistingBashAndGoTargets(t *testing.T) {
	root := t.TempDir()
	mustWriteTestFile(t, filepath.Join(root, "scripts", "present.sh"))
	mustWriteTestFile(t, filepath.Join(root, "scripts", "validator", "main.go"))
	config := parseTestConfig(t, `job:
  script:
    - bash scripts/present.sh
    - go run ./scripts/validator
`)
	if err := validateScriptReferences(root, config); err != nil {
		t.Fatalf("expected existing references to pass, got %v", err)
	}
}

func TestValidateScriptReferencesRejectsUnsafeOrMissingTargets(t *testing.T) {
	// 路径安全覆盖缺失文件、目录穿越和 bash/go 两类绝对路径。
	root := t.TempDir()
	mustWriteTestFile(t, filepath.Join(root, "scripts", "present.sh"))
	outside := filepath.Join(t.TempDir(), "outside.sh")
	mustWriteTestFile(t, outside)

	tests := []struct {
		name    string
		command string
		want    string
	}{
		{name: "missing bash script", command: "bash scripts/missing.sh", want: "does not exist"},
		{name: "traversal", command: "bash scripts/../outside.sh", want: "unsafe"},
		{name: "absolute bash path", command: "bash " + outside, want: "absolute"},
		{name: "absolute go path", command: "go run " + outside, want: "absolute"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := map[string]any{"job": map[string]any{"script": []any{tt.command}}}
			err := validateScriptReferences(root, config)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected %q error, got %v", tt.want, err)
			}
		})
	}
}

func TestValidateScriptReferencesRejectsSymlinkEscape(t *testing.T) {
	// 仓库内名称若通过 symlink 指向仓库外，也必须按真实路径拒绝。
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.sh")
	mustWriteTestFile(t, outside)
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "scripts", "escape.sh")); err != nil {
		t.Fatal(err)
	}
	config := map[string]any{"job": map[string]any{"script": []any{"bash scripts/escape.sh"}}}
	if err := validateScriptReferences(root, config); err == nil || !strings.Contains(err.Error(), "outside repository") {
		t.Fatalf("expected symlink escape error, got %v", err)
	}
}

func TestValidateScriptReferencesChecksBeforeScript(t *testing.T) {
	root := t.TempDir()
	config := map[string]any{"job": map[string]any{"before_script": []any{"bash scripts/missing.sh"}}}
	if err := validateScriptReferences(root, config); err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("expected before_script reference error, got %v", err)
	}
}

func parseTestConfig(t *testing.T, data string) map[string]any {
	t.Helper()
	var config map[string]any
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		t.Fatal(err)
	}
	return config
}

func mustWriteTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("test\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}
