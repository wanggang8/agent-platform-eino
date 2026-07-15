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

SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# canonical 镜像只接受 amd64，避免本机架构悄然改变发布产物。
RUN test "$TARGETARCH" = "amd64"
RUN curl -fsSLo /tmp/go.tar.gz "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    && echo "${GO_SHA256}  /tmp/go.tar.gz" | sha256sum -c - \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz
RUN curl -fsSLo /tmp/node.tar.xz "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz" \
    && echo "${NODE_SHA256}  /tmp/node.tar.xz" | sha256sum -c - \
    && tar -C /usr/local --strip-components=1 -xJf /tmp/node.tar.xz \
    && rm /tmp/node.tar.xz

ENV PATH="/usr/local/go/bin:${PATH}" \
    GOTOOLCHAIN=local \
    PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

RUN test "$(go env GOVERSION)" = "go${GO_VERSION}" \
    && test "$(node --version)" = "v${NODE_VERSION}" \
    && test "$(npm --version)" = "${NPM_VERSION}"
# 浏览器证据直接来自基底镜像，不依赖 clean checkout 尚不存在的 node_modules。
RUN chromium_binary=$(find "/ms-playwright/chromium-${CHROMIUM_REVISION}" \
      -type f -name chrome -perm -111 -print -quit) \
    && test -n "$chromium_binary" \
    && chromium_actual=$("$chromium_binary" --version) \
    && [[ "$chromium_actual" == *"${CHROMIUM_VERSION}"* ]]

WORKDIR /workspace
