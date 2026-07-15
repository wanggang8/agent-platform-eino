# syntax=docker/dockerfile:1.7
ARG PLAYWRIGHT_IMAGE
ARG PLAYWRIGHT_DIGEST
FROM ${PLAYWRIGHT_IMAGE}@${PLAYWRIGHT_DIGEST}

ARG TARGETARCH
ARG GO_VERSION
ARG NODE_VERSION
ARG NPM_VERSION
ARG GO_SHA256
ARG NODE_SHA256
ARG CHROMIUM_REVISION
ARG CHROMIUM_VERSION
ARG UBUNTU_SNAPSHOT
ARG APT_BUILD_PACKAGES

SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# canonical 镜像只接受 amd64，避免本机架构悄然改变发布产物。
RUN test "$TARGETARCH" = "amd64"
# Noble deb822 的每个仓库 stanza 都绑定同一 snapshot，构建依赖不使用浮动 apt 状态。
RUN sources=/etc/apt/sources.list.d/ubuntu.sources \
    && signed_by_count=$(grep -c '^Signed-By:' "$sources") \
    && test "$signed_by_count" -gt 0 \
    && snapshot_before=$(awk '/^Snapshot:/ { count++ } END { print count + 0 }' "$sources") \
    && test "$snapshot_before" = 0 \
    && sed -i "/^Signed-By:/a Snapshot: ${UBUNTU_SNAPSHOT}" "$sources" \
    && snapshot_count=$(grep -Fc "Snapshot: ${UBUNTU_SNAPSHOT}" "$sources") \
    && test "$snapshot_count" = "$signed_by_count" \
    && find /etc/apt/sources.list.d -mindepth 1 -maxdepth 1 -type f ! -name 'ubuntu.sources' -delete \
    && test "$(find /etc/apt/sources.list.d -mindepth 1 -maxdepth 1 -type f -print | wc -l)" = 1 \
    && apt_snapshot_options=(-o Dir::Etc::sourcelist="$sources" -o Dir::Etc::sourceparts="-") \
    && apt-get "${apt_snapshot_options[@]}" update \
    && DEBIAN_FRONTEND=noninteractive apt-get "${apt_snapshot_options[@]}" install -y --no-install-recommends "$APT_BUILD_PACKAGES" \
    && rm -rf /var/lib/apt/lists/*
RUN curl -fsSLo /tmp/go.tar.gz "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    && echo "${GO_SHA256}  /tmp/go.tar.gz" | sha256sum -c - \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz
RUN curl -fsSLo /tmp/node.tar.gz "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.gz" \
    && echo "${NODE_SHA256}  /tmp/node.tar.gz" | sha256sum -c - \
    && tar -C /usr/local --strip-components=1 -xzf /tmp/node.tar.gz \
    && rm /tmp/node.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}" \
    GOTOOLCHAIN=local \
    PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

RUN test "$(go env GOVERSION)" = "go${GO_VERSION}" \
    && test "$(node --version)" = "v${NODE_VERSION}" \
    && test "$(npm --version)" = "${NPM_VERSION}"
# race detector 需要真实 C compiler；在镜像构建时编译并运行最小 race 测试固定该能力。
RUN command -v cc \
    && race_dir=$(mktemp -d) \
    && cd "$race_dir" \
    && go mod init toolchain-race-smoke \
    && printf '%s\n' 'package toolchainrace' '' 'import "testing"' '' \
      'func TestRaceSmoke(t *testing.T) {}' >toolchain_race_test.go \
    && CGO_ENABLED=1 go test -race ./... \
    && cd / \
    && rm -rf "$race_dir"
# 浏览器证据直接来自基底镜像，不依赖 clean checkout 尚不存在的 node_modules。
RUN chromium_binary=$(find "/ms-playwright/chromium-${CHROMIUM_REVISION}" \
      -type f -name chrome -perm -111 -print -quit) \
    && test -n "$chromium_binary" \
    && chromium_actual=$("$chromium_binary" --version) \
    && [[ "$chromium_actual" == *"${CHROMIUM_VERSION}"* ]]

WORKDIR /workspace
