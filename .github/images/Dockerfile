# image-name: ghcr.io/dena/unity-meta-check/unity-meta-check-gh-action
FROM golang:1.16-buster as builder
WORKDIR /go/src/unity-meta-check
COPY . .
RUN make out/gh-action-linux-amd64 && mv ./out/gh-action-linux-amd64 ./out/gh-action

FROM debian:buster-slim
RUN apt-get update \
	&& apt-get install --yes --no-install-recommends git ca-certificates \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/unity-meta-check/out/* /usr/bin/
ENTRYPOINT ["gh-action"]
CMD ["-help"]
