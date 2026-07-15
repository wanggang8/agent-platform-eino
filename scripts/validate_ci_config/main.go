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
	buildJobName          = "toolchain-build"
	verifyJobName         = "toolchain-verify"
	dockerTemplateName    = ".docker-amd64"
	toolchainDotenvPath   = "test-results/toolchain.env"
	baselineScriptPath    = "scripts/run_toolchain_baseline.sh"
	canonicalBuilderPath  = "scripts/build_toolchain_image.sh"
	dockerHostValue       = "tcp://docker:2375"
	dockerBuildkitValue   = "1"
	dockerTLSCertdirValue = ""
	canonicalMkdirLine    = "mkdir -p test-results"
	canonicalLoginLine    = `printf '%s' "$CI_REGISTRY_PASSWORD" | docker login -u "$CI_REGISTRY_USER" --password-stdin "$CI_REGISTRY"`
	canonicalBuilderLine  = `bash ` + canonicalBuilderPath + ` --push "$CI_REGISTRY_IMAGE/toolchain:$CI_COMMIT_SHA" --env-file ` + toolchainDotenvPath
	canonicalGuardLine    = `printf '%s\n' "$TOOLCHAIN_IMAGE" | grep -Eq '^[A-Za-z0-9._/:-]+@sha256:[0-9a-f]{64}$'`
	canonicalPullLine     = `docker pull "$TOOLCHAIN_IMAGE"`
	canonicalBaselineLine = `docker run --rm --platform linux/amd64 -e CI=true -e CI_COMMIT_SHA="$CI_COMMIT_SHA" -v "$CI_PROJECT_DIR:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash ` + baselineScriptPath
)

var (
	lockKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	digestPattern  = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	imagePattern   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/:@-]*$`)
)

// lockValues 仅承载 CI runner 边界需要核对的 Docker CLI 与 DinD 固定值。
type lockValues struct {
	DockerCLIImage   string
	DockerCLIDigest  string
	DockerDindImage  string
	DockerDindDigest string
}

// validateConfig 静态验证 pinned runner、构建到验证的 DAG 与 canonical image 传递。
func validateConfig(data []byte, lock lockValues) error {
	config, err := parseConfig(data)
	if err != nil {
		return err
	}

	stages, ok := stringSlice(config["stages"])
	if !ok || len(stages) != 2 || stages[0] != "toolchain" || stages[1] != "verify" {
		return errors.New("stages must be exactly [toolchain, verify]")
	}

	// 先保留原始未固定 image 的专用诊断，再强制所有有效配置走唯一模板。
	if rawBuild, ok := stringMap(config[buildJobName]); ok {
		if image, _ := rawBuild["image"].(string); image == lock.DockerCLIImage && lock.DockerCLIDigest != "" {
			return errors.New("toolchain-build Docker CLI image must use the lock digest")
		}
	}
	template, ok := stringMap(config[dockerTemplateName])
	if !ok {
		return errors.New("missing .docker-amd64 template")
	}
	build, err := resolveJob(config, template, buildJobName)
	if err != nil {
		return err
	}
	if err := validateDockerBoundary(buildJobName, build, lock); err != nil {
		return err
	}
	if stage, _ := build["stage"].(string); stage != "toolchain" {
		return errors.New("toolchain-build stage must be toolchain")
	}
	buildLines, err := strictScriptLines(build)
	if err != nil || !equalStrings(buildLines, []string{canonicalMkdirLine, canonicalLoginLine, canonicalBuilderLine}) {
		return errors.New("toolchain-build must call scripts/build_toolchain_image.sh with the canonical dotenv path")
	}
	if got := nestedString(build, "artifacts", "reports", "dotenv"); got != toolchainDotenvPath {
		return errors.New("toolchain-build dotenv artifact must be test-results/toolchain.env")
	}

	verify, err := resolveJob(config, template, verifyJobName)
	if err != nil {
		return err
	}
	if err := validateDockerBoundary(verifyJobName, verify, lock); err != nil {
		return err
	}
	if stage, _ := verify["stage"].(string); stage != "verify" {
		return errors.New("toolchain-verify stage must be verify")
	}
	if !hasBuildArtifactNeed(verify) {
		return errors.New("toolchain-verify needs toolchain-build artifacts")
	}

	verifyLines, err := strictScriptLines(verify)
	if err != nil || !equalStrings(verifyLines, []string{canonicalGuardLine, canonicalPullLine, canonicalBaselineLine}) {
		return errors.New("toolchain-verify script must match the canonical @sha256: TOOLCHAIN_IMAGE baseline vector")
	}
	return nil
}

// loadLock 从不可执行的 key=value 文件读取 Docker pins，拒绝重复键与非固定 digest。
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
		if !ok || !lockKeyPattern.MatchString(key) || value == "" || strings.ContainsAny(value, " \t\r\n") {
			return lockValues{}, errors.New("invalid toolchain lock entry")
		}
		if _, exists := values[key]; exists {
			return lockValues{}, fmt.Errorf("duplicate toolchain lock key %s", key)
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return lockValues{}, errors.New("cannot read toolchain lock")
	}

	lock := lockValues{
		DockerCLIImage:   values["DOCKER_CLI_IMAGE"],
		DockerCLIDigest:  values["DOCKER_CLI_AMD64_DIGEST"],
		DockerDindImage:  values["DOCKER_DIND_IMAGE"],
		DockerDindDigest: values["DOCKER_DIND_AMD64_DIGEST"],
	}
	if !imagePattern.MatchString(lock.DockerCLIImage) || !imagePattern.MatchString(lock.DockerDindImage) ||
		!digestPattern.MatchString(lock.DockerCLIDigest) || !digestPattern.MatchString(lock.DockerDindDigest) {
		return lockValues{}, errors.New("toolchain lock must contain pinned Docker CLI and DinD images")
	}
	return lock, nil
}

// validateScriptReferences 确保 CI 中由 bash/go run 启动的仓库脚本真实存在且不能逃逸根目录。
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

func main() {
	configPath := flag.String("config", ".gitlab-ci.yml", "GitLab CI config path")
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
	fmt.Println("GitLab CI config validated")
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

func resolveJob(config, template map[string]any, name string) (map[string]any, error) {
	job, ok := stringMap(config[name])
	if !ok {
		return nil, fmt.Errorf("missing %s job", name)
	}
	extends, ok := job["extends"].(string)
	if !ok || extends != dockerTemplateName || template == nil {
		return nil, fmt.Errorf("%s must extend %s", name, dockerTemplateName)
	}
	resolved := make(map[string]any, len(template)+len(job))
	for key, value := range template {
		resolved[key] = value
	}
	for key, value := range job {
		if key == "variables" {
			inherited, inheritedOK := stringMap(resolved[key])
			local, localOK := stringMap(value)
			if !inheritedOK || !localOK {
				return nil, fmt.Errorf("%s variables must be a map", name)
			}
			resolved[key] = mergeMaps(inherited, local)
			continue
		}
		resolved[key] = value
	}
	return resolved, nil
}

func validateDockerBoundary(jobName string, job map[string]any, lock lockValues) error {
	expectedImage := lock.DockerCLIImage + "@" + lock.DockerCLIDigest
	image, _ := job["image"].(string)
	if image != expectedImage {
		return fmt.Errorf("%s Docker CLI image must use the lock digest", jobName)
	}

	services, ok := anySlice(job["services"])
	if !ok || len(services) != 1 {
		return fmt.Errorf("%s DinD digest must match the lock", jobName)
	}
	service, ok := stringMap(services[0])
	if !ok || service["name"] != lock.DockerDindImage+"@"+lock.DockerDindDigest || service["alias"] != "docker" {
		return fmt.Errorf("%s DinD digest must match the lock", jobName)
	}

	variables, _ := stringMap(job["variables"])
	required := map[string]string{
		"DOCKER_HOST":        dockerHostValue,
		"DOCKER_TLS_CERTDIR": dockerTLSCertdirValue,
		"DOCKER_BUILDKIT":    dockerBuildkitValue,
	}
	for key, value := range required {
		if fmt.Sprint(variables[key]) != value {
			return fmt.Errorf("%s must inherit %s", jobName, key)
		}
	}
	return nil
}

func hasBuildArtifactNeed(job map[string]any) bool {
	needs, ok := anySlice(job["needs"])
	if !ok {
		return false
	}
	for _, item := range needs {
		need, ok := stringMap(item)
		if ok && need["job"] == buildJobName && need["artifacts"] == true {
			return true
		}
	}
	return false
}

func validateRepositoryTarget(root, target string) error {
	if filepath.IsAbs(target) {
		return errors.New("absolute script reference is forbidden")
	}
	target = filepath.ToSlash(target)
	parts := strings.Split(target, "/")
	for _, part := range parts {
		if part == ".." {
			return errors.New("unsafe script reference is forbidden")
		}
	}
	if !strings.HasPrefix(target, "scripts/") && !strings.HasPrefix(target, "./scripts/") {
		return errors.New("script reference must stay under scripts")
	}

	joined := filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(target, "./")))
	realTarget, err := filepath.EvalSymlinks(joined)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("script reference does not exist")
		}
		return errors.New("invalid script reference")
	}
	relative, err := filepath.Rel(root, realTarget)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return errors.New("script reference resolves outside repository")
	}
	return nil
}

func canonicalScriptReferences(config map[string]any) ([]string, error) {
	if containsNonCanonicalScriptKey(config, false) {
		return nil, errors.New("CI config contains a non-canonical script section")
	}
	build, buildOK := stringMap(config[buildJobName])
	verify, verifyOK := stringMap(config[verifyJobName])
	if !buildOK || !verifyOK {
		return nil, errors.New("CI config is missing a canonical job")
	}
	buildLines, buildErr := strictScriptLines(build)
	verifyLines, verifyErr := strictScriptLines(verify)
	if buildErr != nil || verifyErr != nil ||
		!equalStrings(buildLines, []string{canonicalMkdirLine, canonicalLoginLine, canonicalBuilderLine}) ||
		!equalStrings(verifyLines, []string{canonicalGuardLine, canonicalPullLine, canonicalBaselineLine}) {
		return nil, errors.New("CI config contains a non-canonical script vector")
	}
	return []string{canonicalBuilderPath, baselineScriptPath}, nil
}

func containsNonCanonicalScriptKey(value any, insideCanonicalJob bool) bool {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			if key == buildJobName || key == verifyJobName {
				if containsNonCanonicalScriptKey(child, true) {
					return true
				}
				continue
			}
			if key == "script" {
				if !insideCanonicalJob {
					return true
				}
				continue
			}
			if strings.HasSuffix(key, "_script") || containsNonCanonicalScriptKey(child, insideCanonicalJob) {
				return true
			}
		}
	case []any:
		for _, child := range typed {
			if containsNonCanonicalScriptKey(child, insideCanonicalJob) {
				return true
			}
		}
	}
	return false
}

func strictScriptLines(job map[string]any) ([]string, error) {
	items, ok := anySlice(job["script"])
	if !ok {
		return nil, errors.New("script must be a command list")
	}
	lines := make([]string, len(items))
	for index, item := range items {
		line, ok := item.(string)
		if !ok {
			return nil, errors.New("script commands must be strings")
		}
		lines[index] = line
	}
	return lines, nil
}

func equalStrings(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for index := range actual {
		if actual[index] != expected[index] {
			return false
		}
	}
	return true
}

func nestedString(root map[string]any, keys ...string) string {
	var current any = root
	for _, key := range keys {
		mapping, ok := stringMap(current)
		if !ok {
			return ""
		}
		current = mapping[key]
	}
	value, _ := current.(string)
	return value
}

func stringMap(value any) (map[string]any, bool) {
	mapping, ok := value.(map[string]any)
	return mapping, ok
}

func anySlice(value any) ([]any, bool) {
	items, ok := value.([]any)
	return items, ok
}

func stringSlice(value any) ([]string, bool) {
	items, ok := anySlice(value)
	if !ok {
		return nil, false
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		text, ok := item.(string)
		if !ok {
			return nil, false
		}
		result = append(result, text)
	}
	return result, true
}

func mergeMaps(base, override map[string]any) map[string]any {
	merged := make(map[string]any, len(base)+len(override))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range override {
		merged[key] = value
	}
	return merged
}
