# image-name: ghcr.io/dena/unity-meta-check/unity-meta-check-circleci
FROM ghcr.io/dena/unity-meta-check/unity-meta-check-builder:latest as builder

FROM debian:bookworm-slim
# https://circleci.com/docs/2.0/custom-images/#required-tools-for-primary-containers
RUN apt-get update \
	&& apt-get install --yes --no-install-recommends git openssh-server tar gzip ca-certificates \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/unity-meta-check/out/* /usr/bin/
ENTRYPOINT ["unity-meta-check"]
CMD ["-help"]
