package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	canonicalDigestGuardCommand = `printf '%s\n' "$TOOLCHAIN_IMAGE" | grep -Eq '^[A-Za-z0-9._/:-]+@sha256:[0-9a-f]{64}$'`
	canonicalPullCommand        = `docker pull "$TOOLCHAIN_IMAGE"`
	canonicalRunCommand         = `docker run --rm --platform linux/amd64 -e CI=true -e CI_COMMIT_SHA="$CI_COMMIT_SHA" -v "$CI_PROJECT_DIR:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh`
)

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
    - printf '%s' "$CI_REGISTRY_PASSWORD" | docker login -u "$CI_REGISTRY_USER" --password-stdin "$CI_REGISTRY"
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
    - 'printf ''%s\n'' "$TOOLCHAIN_IMAGE" | grep -Eq ''^[A-Za-z0-9._/:-]+@sha256:[0-9a-f]{64}$'''
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
    - mkdir -p test-results
    - printf '%s' "$CI_REGISTRY_PASSWORD" | docker login -u "$CI_REGISTRY_USER" --password-stdin "$CI_REGISTRY"
    - bash scripts/build_toolchain_image.sh --push "$CI_REGISTRY_IMAGE/toolchain:$CI_COMMIT_SHA" --env-file test-results/toolchain.env
  artifacts:
    reports:
      dotenv: test-results/toolchain.env
toolchain-verify:
  extends: .docker-amd64
  stage: verify
  script:
    - 'printf ''%s\n'' "$TOOLCHAIN_IMAGE" | grep -Eq ''^[A-Za-z0-9._/:-]+@sha256:[0-9a-f]{64}$'''
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

func TestValidateRejectsMissingDockerTemplate(t *testing.T) {
	config := parseTestConfig(t, validFixture)
	delete(config, dockerTemplateName)
	err := validateConfig(marshalTestConfig(t, config), completeTestLock())
	if err == nil || !strings.Contains(err.Error(), "template") {
		t.Fatalf("expected missing template error, got %v", err)
	}
}

func TestValidateRejectsNonStringOrMissingExtends(t *testing.T) {
	// 即使 job 内联了与模板等价的 runner 配置，也不能绕过唯一模板继承边界。
	tests := []struct {
		name    string
		jobName string
		value   any
	}{
		{name: "build missing", jobName: buildJobName},
		{name: "verify missing", jobName: verifyJobName},
		{name: "array", jobName: buildJobName, value: []any{dockerTemplateName}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			job := config[tt.jobName].(map[string]any)
			if tt.value == nil {
				delete(job, "extends")
			} else {
				job["extends"] = tt.value
			}
			template := config[dockerTemplateName].(map[string]any)
			job["image"] = template["image"]
			job["services"] = template["services"]
			job["variables"] = template["variables"]
			err := validateConfig(marshalTestConfig(t, config), completeTestLock())
			if err == nil || !strings.Contains(err.Error(), "extend") {
				t.Fatalf("expected exact extends error, got %v", err)
			}
		})
	}
}

func TestValidateAcceptsGitLabStylePartialOverrides(t *testing.T) {
	// variables 采用 map merge；image/services 采用 job 值覆盖模板，最终有效值仍须精确匹配 lock。
	config := parseTestConfig(t, validFixture)
	template := config[dockerTemplateName].(map[string]any)
	build := config[buildJobName].(map[string]any)
	build["variables"] = map[string]any{"DOCKER_BUILDKIT": "1"}
	verify := config[verifyJobName].(map[string]any)
	verify["image"] = template["image"]
	verify["services"] = template["services"]
	if err := validateConfig(marshalTestConfig(t, config), completeTestLock()); err != nil {
		t.Fatalf("expected partial overrides to follow GitLab merge semantics, got %v", err)
	}
}

func TestValidateRejectsInvalidPartialVariableOverride(t *testing.T) {
	config := parseTestConfig(t, validFixture)
	config[buildJobName].(map[string]any)["variables"] = map[string]any{"DOCKER_BUILDKIT": "0"}
	if err := validateConfig(marshalTestConfig(t, config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "DOCKER_BUILDKIT") {
		t.Fatalf("expected invalid local variable override error, got %v", err)
	}
}

func TestValidateRejectsNonMapVariableOverride(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "string", value: "DOCKER_BUILDKIT=1"},
		{name: "array", value: []any{"DOCKER_BUILDKIT=1"}},
		{name: "null", value: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			config[buildJobName].(map[string]any)["variables"] = tt.value
			if err := validateConfig(marshalTestConfig(t, config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "variables") {
				t.Fatalf("expected non-map variables error, got %v", err)
			}
		})
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
	config := strings.Replace(validFixture, "    - 'printf ''%s\\n'' \"$TOOLCHAIN_IMAGE\" | grep -Eq ''^[A-Za-z0-9._/:-]+@sha256:[0-9a-f]{64}$'''\n", "", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "@sha256:") {
		t.Fatalf("expected TOOLCHAIN_IMAGE digest guard error, got %v", err)
	}
}

func TestValidateRejectsDigestTextWithoutGuard(t *testing.T) {
	// 出现 digest 文本不等于运行时校验，必须有针对 TOOLCHAIN_IMAGE 的 guard。
	config := strings.Replace(validFixture, "    - 'printf ''%s\\n'' \"$TOOLCHAIN_IMAGE\" | grep -Eq ''^[A-Za-z0-9._/:-]+@sha256:[0-9a-f]{64}$'''", "    - echo marker-@sha256:value", 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "@sha256:") {
		t.Fatalf("expected non-guard digest text to fail, got %v", err)
	}
}

func TestValidateRejectsNonCanonicalVerifyCommandVector(t *testing.T) {
	// verify 三条必须逐字、定序且唯一，任何 shell 包装或复合命令都改变安全语义。
	tests := []struct {
		name     string
		commands []string
	}{
		{name: "guard comment", commands: []string{"# " + canonicalDigestGuardCommand, canonicalPullCommand, canonicalRunCommand}},
		{name: "guard reversed", commands: []string{"! " + canonicalDigestGuardCommand, canonicalPullCommand, canonicalRunCommand}},
		{name: "guard ignored", commands: []string{canonicalDigestGuardCommand + " || true", canonicalPullCommand, canonicalRunCommand}},
		{name: "guard after pull", commands: []string{canonicalPullCommand, canonicalDigestGuardCommand, canonicalRunCommand}},
		{name: "extra line", commands: []string{canonicalDigestGuardCommand, canonicalPullCommand, "echo extra", canonicalRunCommand}},
		{name: "env wrapper", commands: []string{canonicalDigestGuardCommand, "env " + canonicalPullCommand, canonicalRunCommand}},
		{name: "shell wrapper", commands: []string{canonicalDigestGuardCommand, "sh -c 'docker pull \"$TOOLCHAIN_IMAGE\"'", canonicalRunCommand}},
		{name: "compound command", commands: []string{canonicalDigestGuardCommand, canonicalPullCommand + "; true", canonicalRunCommand}},
		{name: "parenthesized command", commands: []string{canonicalDigestGuardCommand, "(" + canonicalPullCommand + ")", canonicalRunCommand}},
		{name: "extra docker run", commands: []string{canonicalDigestGuardCommand, canonicalPullCommand, "docker run --rm ubuntu:latest true", canonicalRunCommand}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := fixtureWithVerifyCommands(t, tt.commands)
			if err := validateConfig(config, completeTestLock()); err == nil || !strings.Contains(err.Error(), "canonical") {
				t.Fatalf("expected canonical vector error, got %v", err)
			}
		})
	}
}

func fixtureWithVerifyCommands(t *testing.T, commands []string) []byte {
	t.Helper()
	config := parseTestConfig(t, validFixture)
	job := config[verifyJobName].(map[string]any)
	items := make([]any, len(commands))
	for index, command := range commands {
		items[index] = command
	}
	job["script"] = items
	return marshalTestConfig(t, config)
}

func marshalTestConfig(t *testing.T, config map[string]any) []byte {
	t.Helper()
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	return data
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
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "canonical") {
		t.Fatalf("expected additional image run to fail, got %v", err)
	}
}

func TestValidateRejectsToolchainImageUsedAsContainerArgument(t *testing.T) {
	bypass := `docker run --rm ubuntu:latest "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh`
	config := strings.Replace(validFixture, canonicalRunCommand, bypass, 1)
	if err := validateConfig([]byte(config), completeTestLock()); err == nil || !strings.Contains(err.Error(), "canonical") {
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

func TestValidateScriptReferencesAcceptsCanonicalTargets(t *testing.T) {
	root := t.TempDir()
	mustWriteTestFile(t, filepath.Join(root, canonicalBuilderPath))
	mustWriteTestFile(t, filepath.Join(root, baselineScriptPath))
	config := parseTestConfig(t, validFixture)
	if err := validateScriptReferences(root, config); err != nil {
		t.Fatalf("expected canonical references to pass, got %v", err)
	}
}

func TestValidateScriptReferencesRejectsNonCanonicalShellSyntax(t *testing.T) {
	// canonical job 不能通过注释、echo、动态执行、括号或额外 lifecycle hook 引入第二条执行路径。
	root := t.TempDir()
	mustWriteTestFile(t, filepath.Join(root, canonicalBuilderPath))
	mustWriteTestFile(t, filepath.Join(root, baselineScriptPath))
	mustWriteTestFile(t, filepath.Join(root, "scripts", "validator", "main.go"))
	canonicalBuild := []string{canonicalMkdirLine, canonicalLoginLine, canonicalBuilderLine}
	tests := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "echo", mutate: func(config map[string]any) {
			setJobCommands(config, buildJobName, []string{canonicalMkdirLine, canonicalLoginLine, "echo " + canonicalBuilderLine})
		}},
		{name: "command substitution", mutate: func(config map[string]any) {
			setJobCommands(config, buildJobName, []string{canonicalMkdirLine, canonicalLoginLine, "$(" + canonicalBuilderLine + ")"})
		}},
		{name: "parentheses", mutate: func(config map[string]any) {
			setJobCommands(config, buildJobName, []string{canonicalMkdirLine, canonicalLoginLine, "(" + canonicalBuilderLine + ")"})
		}},
		{name: "extra go run", mutate: func(config map[string]any) {
			setJobCommands(config, buildJobName, append(append([]string{}, canonicalBuild...), "go run ./scripts/validator"))
		}},
		{name: "before script", mutate: func(config map[string]any) { config["before_script"] = []any{"echo before"} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseTestConfig(t, validFixture)
			tt.mutate(config)
			err := validateScriptReferences(root, config)
			if err == nil || !strings.Contains(err.Error(), "canonical") {
				t.Fatalf("expected canonical script error, got %v", err)
			}
		})
	}
}

func TestValidateScriptReferencesRejectsMissingCanonicalTarget(t *testing.T) {
	root := t.TempDir()
	mustWriteTestFile(t, filepath.Join(root, canonicalBuilderPath))
	if err := validateScriptReferences(root, parseTestConfig(t, validFixture)); err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("expected missing baseline error, got %v", err)
	}
}

func TestValidateScriptReferencesRejectsCanonicalSymlinkEscape(t *testing.T) {
	// 仓库内名称若通过 symlink 指向仓库外，也必须按真实路径拒绝。
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.sh")
	mustWriteTestFile(t, outside)
	mustWriteTestFile(t, filepath.Join(root, canonicalBuilderPath))
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, baselineScriptPath)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, baselineScriptPath)); err != nil {
		t.Fatal(err)
	}
	config := parseTestConfig(t, validFixture)
	if err := validateScriptReferences(root, config); err == nil || !strings.Contains(err.Error(), "outside repository") {
		t.Fatalf("expected symlink escape error, got %v", err)
	}
}

func TestValidateRepositoryTargetRejectsUnsafePaths(t *testing.T) {
	// 严格 vector 之外仍保留底层路径边界，直接证明绝对路径与目录穿越被拒绝。
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.sh")
	mustWriteTestFile(t, outside)
	tests := []struct {
		name   string
		target string
		want   string
	}{
		{name: "traversal", target: "scripts/../outside.sh", want: "unsafe"},
		{name: "absolute", target: outside, want: "absolute"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRepositoryTarget(root, tt.target); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected %s path error, got %v", tt.want, err)
			}
		})
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

func setJobCommands(config map[string]any, jobName string, commands []string) {
	job := config[jobName].(map[string]any)
	items := make([]any, len(commands))
	for index, command := range commands {
		items[index] = command
	}
	job["script"] = items
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
