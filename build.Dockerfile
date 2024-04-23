# image-name: ghcr.io/dena/unity-meta-check/unity-meta-check-builder
FROM golang:1.22.2-bookworm as builder
ARG TARGETARCH
WORKDIR /go/src/unity-meta-check
COPY . .
RUN make -j$(nproc) out/unity-meta-check-linux-${TARGETARCH} \
				out/unity-meta-check-junit-linux-${TARGETARCH} \
				out/unity-meta-check-github-pr-comment-linux-${TARGETARCH} \
				out/unity-meta-autofix-linux-${TARGETARCH} \
				out/gh-action-linux-${TARGETARCH} && \
			mv ./out/unity-meta-check-linux-${TARGETARCH} ./out/unity-meta-check && \
			mv ./out/unity-meta-check-junit-linux-${TARGETARCH} ./out/unity-meta-check-junit && \
			mv ./out/unity-meta-check-github-pr-comment-linux-${TARGETARCH} ./out/unity-meta-check-github-pr-comment && \
			mv ./out/unity-meta-autofix-linux-${TARGETARCH} ./out/unity-meta-autofix && \
			mv ./out/gh-action-linux-${TARGETARCH} ./out/gh-action
