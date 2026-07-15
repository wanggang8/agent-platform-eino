package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	buildJobName  = "toolchain-build"
	verifyJobName = "toolchain-verify"

	canonicalCredentialConfigCommand = `install -d -m 0700 "$RUNNER_TEMP/docker-config" "$RUNNER_TEMP/docker-credential-bin"
printf '%s\n' '{"credHelpers":{"ghcr.io":"github-token"}}' > "$RUNNER_TEMP/docker-config/config.json"
install -m 0700 scripts/docker-credential-github-token "$RUNNER_TEMP/docker-credential-bin/docker-credential-github-token"
printf 'DOCKER_CONFIG=%s\n' "$RUNNER_TEMP/docker-config" >> "$GITHUB_ENV"
printf '%s\n' "$RUNNER_TEMP/docker-credential-bin" >> "$GITHUB_PATH"`
	canonicalBuildCommand = `repository=${GITHUB_REPOSITORY,,}
image_ref="ghcr.io/$repository/toolchain:$GITHUB_SHA"
bash scripts/build_toolchain_image.sh --push "$image_ref" --env-file test-results/toolchain.env
cat test-results/toolchain.env >> "$GITHUB_OUTPUT"`
	canonicalVerifyCommand = `[[ "$TOOLCHAIN_IMAGE" =~ ^ghcr\.io/[a-z0-9._/-]+/toolchain:[0-9a-f]{40}@sha256:[0-9a-f]{64}$ ]]
docker pull "$TOOLCHAIN_IMAGE"
docker run --rm --platform linux/amd64 -e CI=true -e GITHUB_SHA="$GITHUB_SHA" -v "$GITHUB_WORKSPACE:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh`

	canonicalBuildOutput  = "${{ steps.build.outputs.TOOLCHAIN_IMAGE }}"
	canonicalVerifyImage  = "${{ needs.toolchain-build.outputs.toolchain_image }}"
	canonicalGitHubToken  = "${{ secrets.GITHUB_TOKEN }}"
	canonicalGitHubActor  = "${{ github.actor }}"
	canonicalAlways       = "${{ always() }}"
	canonicalArtifactPath = "test-results/toolchain-baseline.log\ntest-results/eino-workbench-playwright-report/\n"
)

var (
	lockKeyPattern       = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	lockValuePattern     = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/:@+-]*$`)
	actionSHAPattern     = regexp.MustCompile(`^[0-9a-f]{40}$`)
	digestPattern        = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	exactSemverPattern   = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
	buildKitImagePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._/-]*:v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
	gitLabSyntaxPattern  = regexp.MustCompile(`(?i)\.gitlab-ci|\bCI_(?:COMMIT|REGISTRY|PROJECT|PIPELINE|JOB|SERVER|RUNNER|DEFAULT_BRANCH|MERGE_REQUEST)[A-Z0-9_]*\b`)
	scriptCallPattern    = regexp.MustCompile(`(?:^|[[:space:]])(?:bash|go[[:space:]]+run)[[:space:]]+([^[:space:]"']+)`)
	helperInstallPattern = regexp.MustCompile(`(?:^|[[:space:]])install[[:space:]]+-m[[:space:]]+0700[[:space:]]+(scripts/docker-credential-[^[:space:]"']+)`)
	toolchainLockKeys    = []string{
		"PLATFORM",
		"PLAYWRIGHT_IMAGE",
		"PLAYWRIGHT_AMD64_DIGEST",
		"PLAYWRIGHT_VERSION",
		"CHROMIUM_REVISION",
		"CHROMIUM_VERSION",
		"BASE_OS",
		"FONT_POLICY",
		"UBUNTU_SNAPSHOT",
		"APT_BUILD_PACKAGES",
		"GO_LINUX_AMD64_SHA256",
		"NODE_LINUX_X64_SHA256",
		"GITHUB_RUNNER",
		"ACTIONS_CHECKOUT_SHA",
		"ACTIONS_UPLOAD_ARTIFACT_SHA",
		"DOCKER_SETUP_DOCKER_SHA",
		"DOCKER_SETUP_BUILDX_SHA",
		"DOCKER_ENGINE_VERSION",
		"DOCKER_BUILDX_VERSION",
		"BUILDKIT_IMAGE",
		"BUILDKIT_DIGEST",
	}
	toolchainLockKeySet = func() map[string]struct{} {
		keys := make(map[string]struct{}, len(toolchainLockKeys))
		for _, key := range toolchainLockKeys {
			keys[key] = struct{}{}
		}
		return keys
	}()
)

// lockValues 仅承载 GitHub Actions 执行边界需要核对的固定输入。
type lockValues struct {
	GitHubRunner             string
	ActionsCheckoutSHA       string
	ActionsUploadArtifactSHA string
	DockerSetupDockerSHA     string
	DockerSetupBuildxSHA     string
	DockerEngineVersion      string
	DockerBuildxVersion      string
	BuildKitImage            string
	BuildKitDigest           string
}

// validateConfig 对 GitHub Actions workflow 的触发器、权限、DAG 与执行入口做精确静态校验。
func validateConfig(data []byte, lock lockValues) error {
	if gitLabSyntaxPattern.Match(data) {
		return errors.New("GitLab syntax is forbidden")
	}
	config, err := parseConfig(data)
	if err != nil {
		return err
	}
	if !hasExactKeys(config, "name", "on", "permissions", "jobs") {
		return errors.New("workflow top-level keys must be exact")
	}
	if config["name"] != "Toolchain Gate" {
		return errors.New("workflow name must be Toolchain Gate")
	}
	if err := validateTriggers(config); err != nil {
		return err
	}
	if err := validateTopLevelPermissions(config); err != nil {
		return err
	}
	jobs, ok := stringMap(config["jobs"])
	if !ok || !hasExactKeys(jobs, buildJobName, verifyJobName) {
		return errors.New("jobs must be exactly toolchain-build and toolchain-verify")
	}
	build, buildOK := stringMap(jobs[buildJobName])
	verify, verifyOK := stringMap(jobs[verifyJobName])
	if !buildOK || !verifyOK {
		return errors.New("jobs must be exactly toolchain-build and toolchain-verify")
	}
	if err := validateBuildJob(build, lock); err != nil {
		return err
	}
	return validateVerifyJob(verify, lock)
}

// loadLock 读取不可执行的 KEY=VALUE lock，并只投影 GitHub Actions 所需的九个固定字段。
func loadLock(path string) (lockValues, error) {
	file, err := os.Open(path)
	if err != nil {
		return lockValues{}, errors.New("cannot read toolchain lock")
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || !lockKeyPattern.MatchString(key) || !lockValuePattern.MatchString(value) {
			return lockValues{}, errors.New("invalid toolchain lock entry")
		}
		if _, known := toolchainLockKeySet[key]; !known {
			return lockValues{}, fmt.Errorf("unknown toolchain lock key %s", key)
		}
		if _, exists := values[key]; exists {
			return lockValues{}, fmt.Errorf("duplicate toolchain lock key %s", key)
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return lockValues{}, errors.New("cannot read toolchain lock")
	}

	// 完整 schema 保持 lock 的封闭性；返回值仍只投影 workflow 使用的九个字段。
	for _, key := range toolchainLockKeys {
		if _, ok := values[key]; !ok {
			return lockValues{}, fmt.Errorf("missing toolchain lock key %s", key)
		}
	}

	if values["GITHUB_RUNNER"] != "ubuntu-24.04" {
		return lockValues{}, errors.New("invalid GITHUB_RUNNER pin")
	}
	for _, key := range []string{"ACTIONS_CHECKOUT_SHA", "ACTIONS_UPLOAD_ARTIFACT_SHA", "DOCKER_SETUP_DOCKER_SHA", "DOCKER_SETUP_BUILDX_SHA"} {
		if !actionSHAPattern.MatchString(values[key]) {
			return lockValues{}, fmt.Errorf("invalid %s pin", key)
		}
	}
	for _, key := range []string{"DOCKER_ENGINE_VERSION", "DOCKER_BUILDX_VERSION"} {
		if !exactSemverPattern.MatchString(values[key]) {
			return lockValues{}, fmt.Errorf("invalid %s pin", key)
		}
	}
	if !buildKitImagePattern.MatchString(values["BUILDKIT_IMAGE"]) {
		return lockValues{}, errors.New("invalid BUILDKIT_IMAGE pin")
	}
	if !digestPattern.MatchString(values["BUILDKIT_DIGEST"]) {
		return lockValues{}, errors.New("invalid BUILDKIT_DIGEST pin")
	}

	return lockValues{
		GitHubRunner:             values["GITHUB_RUNNER"],
		ActionsCheckoutSHA:       values["ACTIONS_CHECKOUT_SHA"],
		ActionsUploadArtifactSHA: values["ACTIONS_UPLOAD_ARTIFACT_SHA"],
		DockerSetupDockerSHA:     values["DOCKER_SETUP_DOCKER_SHA"],
		DockerSetupBuildxSHA:     values["DOCKER_SETUP_BUILDX_SHA"],
		DockerEngineVersion:      values["DOCKER_ENGINE_VERSION"],
		DockerBuildxVersion:      values["DOCKER_BUILDX_VERSION"],
		BuildKitImage:            values["BUILDKIT_IMAGE"],
		BuildKitDigest:           values["BUILDKIT_DIGEST"],
	}, nil
}

// validateTriggers 只允许 push 与手动触发，避免未评审的执行入口。
func validateTriggers(config map[string]any) error {
	triggers, ok := stringMap(config["on"])
	if !ok || !hasExactKeys(triggers, "push", "workflow_dispatch") || triggers["push"] != nil || triggers["workflow_dispatch"] != nil {
		return errors.New("workflow triggers must be exactly push and workflow_dispatch")
	}
	return nil
}

// validateTopLevelPermissions 要求默认 token 权限为空，权限只能由两个 job 显式授予。
func validateTopLevelPermissions(config map[string]any) error {
	permissions, ok := stringMap(config["permissions"])
	if !ok || len(permissions) != 0 {
		return errors.New("top-level permissions must be empty")
	}
	return nil
}

// validateBuildJob 固定构建 job 的 runner、最小权限、输出与五个有序步骤。
func validateBuildJob(job map[string]any, lock lockValues) error {
	if job["runs-on"] != lock.GitHubRunner {
		return errors.New("toolchain-build runner must match lock")
	}
	if err := validateExactPermissions(buildJobName, job["permissions"], map[string]string{"contents": "read", "packages": "write"}); err != nil {
		return err
	}
	outputs, ok := stringMap(job["outputs"])
	if !ok || !hasExactKeys(outputs, "toolchain_image") || outputs["toolchain_image"] != canonicalBuildOutput {
		return errors.New("toolchain-build output must use steps.build.outputs.TOOLCHAIN_IMAGE")
	}
	steps, ok := anySlice(job["steps"])
	if !ok || len(steps) != 5 {
		return errors.New("toolchain-build steps must be exact")
	}
	checkout, checkoutOK := stringMap(steps[0])
	setupDocker, setupDockerOK := stringMap(steps[1])
	configure, configureOK := stringMap(steps[2])
	setupBuildx, setupBuildxOK := stringMap(steps[3])
	build, buildOK := stringMap(steps[4])
	if !checkoutOK || !setupDockerOK || !configureOK || !setupBuildxOK || !buildOK {
		return errors.New("toolchain-build steps must be exact")
	}
	if err := validateActionStep(checkout, "actions/checkout@"+lock.ActionsCheckoutSHA, map[string]string{"fetch-depth": "0"}); err != nil {
		return err
	}
	if err := validateActionStep(setupDocker, "docker/setup-docker-action@"+lock.DockerSetupDockerSHA, map[string]string{"version": "v" + lock.DockerEngineVersion}); err != nil {
		return err
	}
	if err := validateCredentialConfigStep(configure); err != nil {
		return err
	}
	if err := validateActionStep(setupBuildx, "docker/setup-buildx-action@"+lock.DockerSetupBuildxSHA, map[string]string{
		"version":     "v" + lock.DockerBuildxVersion,
		"driver-opts": "image=" + lock.BuildKitImage + "@" + lock.BuildKitDigest,
	}); err != nil {
		return err
	}
	if !hasExactKeys(build, "name", "id", "env", "run") || build["name"] != "Build and push canonical image" || validateRunStep(build, "build", canonicalBuildCommand) != nil {
		return errors.New("Build and push canonical image step must match canonical command")
	}
	if !exactStringMap(build["env"], canonicalCredentialEnv()) {
		return errors.New("Build and push canonical image authentication must be step-local")
	}
	if !hasExactKeys(job, "runs-on", "permissions", "outputs", "steps") {
		return errors.New("toolchain-build keys must be exact")
	}
	return nil
}

// validateVerifyJob 固定验证 job 的依赖、只读权限、镜像来源与五个有序步骤。
func validateVerifyJob(job map[string]any, lock lockValues) error {
	if job["needs"] != buildJobName {
		return errors.New("toolchain-verify must need toolchain-build")
	}
	if job["runs-on"] != lock.GitHubRunner {
		return errors.New("toolchain-verify runner must match lock")
	}
	if err := validateExactPermissions(verifyJobName, job["permissions"], map[string]string{"contents": "read", "packages": "read"}); err != nil {
		return err
	}
	if !exactStringMap(job["env"], map[string]string{"TOOLCHAIN_IMAGE": canonicalVerifyImage}) {
		return errors.New("toolchain-verify image must come from toolchain-build output")
	}
	steps, ok := anySlice(job["steps"])
	if !ok || len(steps) != 5 {
		return errors.New("toolchain-verify steps must be exact")
	}
	checkout, checkoutOK := stringMap(steps[0])
	setupDocker, setupDockerOK := stringMap(steps[1])
	configure, configureOK := stringMap(steps[2])
	verify, verifyOK := stringMap(steps[3])
	artifact, artifactOK := stringMap(steps[4])
	if !checkoutOK || !setupDockerOK || !configureOK || !verifyOK || !artifactOK {
		return errors.New("toolchain-verify steps must be exact")
	}
	if err := validateActionStep(checkout, "actions/checkout@"+lock.ActionsCheckoutSHA, map[string]string{"fetch-depth": "0"}); err != nil {
		return err
	}
	if err := validateActionStep(setupDocker, "docker/setup-docker-action@"+lock.DockerSetupDockerSHA, map[string]string{"version": "v" + lock.DockerEngineVersion}); err != nil {
		return err
	}
	if err := validateCredentialConfigStep(configure); err != nil {
		return err
	}
	if !hasExactKeys(verify, "name", "env", "run") || verify["name"] != "Verify digest image" || validateRunStep(verify, "", canonicalVerifyCommand) != nil {
		return errors.New("Verify digest image step must match canonical command")
	}
	if !exactStringMap(verify["env"], canonicalCredentialEnv()) {
		return errors.New("Verify digest image authentication must be step-local")
	}
	if err := validateArtifactStep(artifact, lock); err != nil {
		return err
	}
	if !hasExactKeys(job, "needs", "runs-on", "permissions", "env", "steps") {
		return errors.New("toolchain-verify keys must be exact")
	}
	return nil
}

// validateCredentialConfigStep 固定无 token 的 Docker 配置与 helper 安装位置。
func validateCredentialConfigStep(step map[string]any) error {
	if !hasExactKeys(step, "name", "run") ||
		step["name"] != "Configure memory-only GHCR credentials" ||
		validateRunStep(step, "", canonicalCredentialConfigCommand) != nil {
		return errors.New("credential helper configuration must match canonical command")
	}
	return nil
}

// canonicalCredentialEnv 只允许真正访问 GHCR 的 step 获取当前进程凭据。
func canonicalCredentialEnv() map[string]string {
	return map[string]string{
		"GHCR_ACTOR": canonicalGitHubActor,
		"GHCR_TOKEN": canonicalGitHubToken,
	}
}

// validateExactPermissions 拒绝缺失、升级或额外的 job token 权限。
func validateExactPermissions(jobName string, actual any, expected map[string]string) error {
	if !exactStringMap(actual, expected) {
		return fmt.Errorf("%s permissions must be exact", jobName)
	}
	return nil
}

// validateActionStep 固定 action 名称、commit SHA、输入集合和值，禁止 tag、branch 或额外参数。
func validateActionStep(step map[string]any, expectedUses string, expectedWith map[string]string) error {
	name, summary := actionIdentity(expectedUses)
	if !hasExactKeys(step, "name", "uses", "with") || step["name"] != name || step["uses"] != expectedUses || !exactStringMap(step["with"], expectedWith) {
		return errors.New(summary)
	}
	return nil
}

// validateRunStep 对 run 文本去除 YAML 外围空白后逐字比较，并固定 id 是否存在。
func validateRunStep(step map[string]any, expectedID, expectedRun string) error {
	run, ok := step["run"].(string)
	if !ok || strings.TrimSpace(run) != strings.TrimSpace(expectedRun) {
		return errors.New("run step command mismatch")
	}
	if expectedID == "" {
		if _, exists := step["id"]; exists {
			return errors.New("run step id mismatch")
		}
		return nil
	}
	if step["id"] != expectedID {
		return errors.New("run step id mismatch")
	}
	return nil
}

// validateArtifactStep 固定失败也上传的证据集合、保留期与 action commit。
func validateArtifactStep(step map[string]any, lock lockValues) error {
	expectedWith := map[string]string{
		"name":              "toolchain-evidence-${{ github.sha }}",
		"path":              canonicalArtifactPath,
		"if-no-files-found": "error",
		"retention-days":    "30",
	}
	if !hasExactKeys(step, "name", "if", "uses", "with") ||
		step["name"] != "Upload toolchain evidence" ||
		step["if"] != canonicalAlways ||
		step["uses"] != "actions/upload-artifact@"+lock.ActionsUploadArtifactSHA ||
		!exactStringMap(step["with"], expectedWith) {
		return errors.New("artifact step must match lock")
	}
	return nil
}

// canonicalScriptReferences 只解析 workflow run 步骤中实际交给 bash 或 go run 的仓库路径。
func canonicalScriptReferences(config map[string]any) ([]string, error) {
	jobs, ok := stringMap(config["jobs"])
	if !ok {
		return nil, errors.New("jobs must be a map")
	}
	seen := make(map[string]struct{})
	references := make([]string, 0, 2)
	for _, jobName := range []string{buildJobName, verifyJobName} {
		job, ok := stringMap(jobs[jobName])
		if !ok {
			return nil, fmt.Errorf("missing %s job", jobName)
		}
		steps, ok := anySlice(job["steps"])
		if !ok {
			return nil, fmt.Errorf("%s steps must be a list", jobName)
		}
		for _, rawStep := range steps {
			step, ok := stringMap(rawStep)
			if !ok {
				return nil, fmt.Errorf("%s step must be a map", jobName)
			}
			run, _ := step["run"].(string)
			for _, match := range scriptCallPattern.FindAllStringSubmatch(run, -1) {
				target := match[1]
				if _, exists := seen[target]; exists {
					continue
				}
				seen[target] = struct{}{}
				references = append(references, target)
			}
			for _, match := range helperInstallPattern.FindAllStringSubmatch(run, -1) {
				target := match[1]
				if _, exists := seen[target]; exists {
					continue
				}
				seen[target] = struct{}{}
				references = append(references, target)
			}
		}
	}
	return references, nil
}

// validateScriptReferences 对解析到的仓库脚本执行真实路径检查，阻断绝对路径、父级穿越与 symlink 逃逸。
func validateScriptReferences(root string, config map[string]any) error {
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return errors.New("invalid repository root")
	}
	references, err := canonicalScriptReferences(config)
	if err != nil {
		return err
	}
	for _, target := range references {
		if err := validateRepositoryTarget(realRoot, target); err != nil {
			return err
		}
	}
	return nil
}

func validateRepositoryTarget(realRoot, target string) error {
	if filepath.IsAbs(target) {
		return errors.New("script reference escapes repository")
	}
	cleanTarget := filepath.Clean(target)
	if cleanTarget == ".." || strings.HasPrefix(cleanTarget, ".."+string(filepath.Separator)) || !strings.HasPrefix(filepath.ToSlash(cleanTarget), "scripts/") {
		return errors.New("script reference escapes repository")
	}
	realTarget, err := filepath.EvalSymlinks(filepath.Join(realRoot, cleanTarget))
	if err != nil {
		return errors.New("script reference does not exist")
	}
	relative, err := filepath.Rel(realRoot, realTarget)
	if err != nil || filepath.IsAbs(relative) || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || !strings.HasPrefix(filepath.ToSlash(relative), "scripts/") {
		return errors.New("script reference escapes repository")
	}
	info, err := os.Stat(realTarget)
	if err != nil || !info.Mode().IsRegular() {
		return errors.New("script reference does not exist")
	}
	return nil
}

func main() {
	configPath := flag.String("config", ".github/workflows/toolchain.yml", "GitHub Actions workflow path")
	lockPath := flag.String("lock", "build/toolchain/toolchain.lock", "toolchain lock path")
	root := flag.String("root", ".", "repository root")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		exitInvalid(errors.New("cannot read CI config"))
	}
	lock, err := loadLock(*lockPath)
	if err != nil {
		exitInvalid(err)
	}
	if err := validateConfig(data, lock); err != nil {
		exitInvalid(err)
	}
	config, err := parseConfig(data)
	if err != nil {
		exitInvalid(err)
	}
	if err := validateScriptReferences(*root, config); err != nil {
		exitInvalid(err)
	}
	fmt.Println("GitHub Actions config validated")
}

func exitInvalid(err error) {
	fmt.Fprintf(os.Stderr, "CI config invalid: %s\n", err)
	os.Exit(1)
}

func parseConfig(data []byte) (map[string]any, error) {
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.New("invalid YAML")
	}
	if config == nil {
		return nil, errors.New("invalid YAML document")
	}
	return config, nil
}

func actionIdentity(expectedUses string) (string, string) {
	switch {
	case strings.HasPrefix(expectedUses, "actions/checkout@"):
		return "Checkout", "Checkout action must use the lock SHA"
	case strings.HasPrefix(expectedUses, "docker/setup-docker-action@"):
		return "Setup Docker", "Setup Docker action must match lock"
	case strings.HasPrefix(expectedUses, "docker/setup-buildx-action@"):
		return "Setup Buildx", "Setup Buildx action must match lock"
	default:
		return "", "action step must match lock"
	}
}

func exactStringMap(actual any, expected map[string]string) bool {
	values, ok := stringMap(actual)
	if !ok || len(values) != len(expected) {
		return false
	}
	for key, expectedValue := range expected {
		value, ok := scalarString(values[key])
		if !ok || value != expectedValue {
			return false
		}
	}
	return true
}

func scalarString(value any) (string, bool) {
	switch typed := value.(type) {
	case string:
		return typed, true
	case int:
		return fmt.Sprintf("%d", typed), true
	case int64:
		return fmt.Sprintf("%d", typed), true
	default:
		return "", false
	}
}

func hasExactKeys(values map[string]any, expected ...string) bool {
	if len(values) != len(expected) {
		return false
	}
	for _, key := range expected {
		if _, ok := values[key]; !ok {
			return false
		}
	}
	return true
}

func stringMap(value any) (map[string]any, bool) {
	values, ok := value.(map[string]any)
	return values, ok
}

func anySlice(value any) ([]any, bool) {
	values, ok := value.([]any)
	return values, ok
}
