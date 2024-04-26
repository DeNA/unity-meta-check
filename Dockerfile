# image-name: ghcr.io/dena/unity-meta-check/unity-meta-check
FROM ghcr.io/dena/unity-meta-check/unity-meta-check-builder:latest as builder

FROM debian:bookworm-slim
RUN apt-get update \
	&& apt-get install --yes --no-install-recommends git \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/unity-meta-check/out/* /usr/bin/
ENTRYPOINT ["unity-meta-check"]
CMD ["-help"]
