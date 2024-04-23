# image-name: ghcr.io/dena/unity-meta-check/unity-meta-check-circleci
FROM golang:1.22-bookworm as builder
ARG TARGETARCH
WORKDIR /go/src/unity-meta-check
COPY . .
RUN make out/unity-meta-check-linux-${TARGETARCH} out/unity-meta-check-junit-linux-${TARGETARCH} out/unity-meta-check-github-pr-comment-linux-${TARGETARCH} out/unity-meta-autofix-linux-${TARGETARCH} && \
	mv ./out/unity-meta-check-linux-${TARGETARCH} ./out/unity-meta-check && \
	mv ./out/unity-meta-check-junit-linux-${TARGETARCH} ./out/unity-meta-check-junit && \
	mv ./out/unity-meta-check-github-pr-comment-linux-${TARGETARCH} ./out/unity-meta-check-github-pr-comment && \
	mv ./out/unity-meta-autofix-linux-${TARGETARCH} ./out/unity-meta-autofix

FROM debian:bookworm-slim
# https://circleci.com/docs/2.0/custom-images/#required-tools-for-primary-containers
RUN apt-get update \
	&& apt-get install --yes --no-install-recommends git openssh-server tar gzip ca-certificates \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/unity-meta-check/out/* /usr/bin/
ENTRYPOINT ["unity-meta-check"]
CMD ["-help"]
