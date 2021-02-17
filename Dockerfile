FROM golang:1.15-alpine3.12 as builder
RUN apk add --no-cache make
WORKDIR /go/src/unity-meta-check
COPY . .
RUN make linux-amd64 && \
	mv ./dist/unity-meta-check-linux-amd64 ./dist/unity-meta-check && \
	mv ./dist/unity-meta-check-junit-linux-amd64 ./dist/unity-meta-check-junit && \
	mv ./dist/unity-meta-check-github-pr-comment-linux-amd64 ./dist/unity-meta-check-github-pr-comment && \
	mv ./dist/unity-meta-autofix-linux-amd64 ./dist/unity-meta-autofix

FROM debian:buster-slim
# https://circleci.com/docs/2.0/custom-images/#required-tools-for-primary-containers
RUN apt-get update && apt-get install --no-install-recommends --yes git openssh-server tar gzip ca-certificates
COPY --from=builder /go/src/unity-meta-check/dist/* /usr/bin/
ENTRYPOINT ["unity-meta-check"]
CMD ["-help"]
