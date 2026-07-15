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
	canonicalBuilderLine  = `bash ` + canonicalBuilderPath + ` --push "$CI_REGISTRY_IMAGE/toolchain:$CI_COMMIT_SHA" --env-file ` + toolchainDotenvPath
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

	template, _ := stringMap(config[dockerTemplateName])
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
	if !hasCanonicalBuilderInvocation(scriptLines(build)) {
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

	if err := validateVerifyCommands(scriptLines(verify)); err != nil {
		return err
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
	for _, command := range allScriptLines(config) {
		fields := strings.Fields(command)
		for index := 0; index < len(fields); index++ {
			var target string
			switch {
			case fields[index] == "bash" && index+1 < len(fields):
				target = trimShellToken(fields[index+1])
			case fields[index] == "go" && index+2 < len(fields) && fields[index+1] == "run":
				target = trimShellToken(fields[index+2])
			default:
				continue
			}
			if err := validateRepositoryTarget(realRoot, target); err != nil {
				return err
			}
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
	resolved := make(map[string]any, len(template)+len(job))
	if extends, _ := job["extends"].(string); extends != "" {
		if extends != dockerTemplateName || template == nil {
			return nil, fmt.Errorf("%s must extend %s", name, dockerTemplateName)
		}
		for key, value := range template {
			resolved[key] = value
		}
	}
	for key, value := range job {
		if key == "variables" {
			inherited, _ := stringMap(resolved[key])
			local, _ := stringMap(value)
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

func hasCanonicalBuilderInvocation(lines []string) bool {
	for _, line := range lines {
		if strings.TrimSpace(line) == canonicalBuilderLine {
			return true
		}
	}
	return false
}

func validateVerifyCommands(lines []string) error {
	guardFound := false
	pullFound := false
	runFound := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, `case "$TOOLCHAIN_IMAGE" in *@sha256:*)`) && strings.Contains(line, "exit 1") {
			guardFound = true
		}
		if strings.HasPrefix(line, "docker pull ") {
			if line != `docker pull "$TOOLCHAIN_IMAGE"` || pullFound {
				return errors.New("toolchain-verify must pull only TOOLCHAIN_IMAGE")
			}
			pullFound = true
		}
		if strings.HasPrefix(line, "docker run ") {
			if runFound || line != canonicalBaselineLine {
				return errors.New("toolchain-verify may run only TOOLCHAIN_IMAGE with the canonical baseline")
			}
			runFound = true
		}
	}
	if !guardFound {
		return errors.New("toolchain-verify must reject TOOLCHAIN_IMAGE without @sha256:")
	}
	if !pullFound {
		return errors.New("toolchain-verify must pull only TOOLCHAIN_IMAGE")
	}
	if !runFound {
		return errors.New("toolchain-verify must call the canonical baseline")
	}
	return nil
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

func allScriptLines(value any) []string {
	var lines []string
	var walk func(any)
	walk = func(current any) {
		switch typed := current.(type) {
		case map[string]any:
			for key, child := range typed {
				if key == "script" || strings.HasSuffix(key, "_script") {
					lines = append(lines, scriptValueLines(child)...)
				}
				walk(child)
			}
		case []any:
			for _, child := range typed {
				walk(child)
			}
		}
	}
	walk(value)
	return lines
}

func scriptLines(job map[string]any) []string {
	return scriptValueLines(job["script"])
}

func scriptValueLines(value any) []string {
	items, ok := anySlice(value)
	if !ok {
		if line, ok := value.(string); ok {
			return []string{line}
		}
		return nil
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		if line, ok := item.(string); ok {
			lines = append(lines, line)
		}
	}
	return lines
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

func trimShellToken(value string) string {
	return strings.Trim(value, `"'`)
}
