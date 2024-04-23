# image-name: ghcr.io/dena/unity-meta-check/unity-meta-check-gh-action
FROM ghcr.io/dena/unity-meta-check/unity-meta-check-builder:latest as builder

FROM debian:bookworm-slim
RUN apt-get update \
	&& apt-get install --yes --no-install-recommends git ca-certificates \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/unity-meta-check/out/gh-action /usr/bin/
ENTRYPOINT ["gh-action"]
CMD ["-help"]
