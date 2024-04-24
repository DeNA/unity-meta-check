all: darwin-amd64 linux-amd64 windows-amd64 darwin-arm64 linux-arm64 windows-arm64


clean:
	go clean
	git clean -fdx ./out


test:
	go test ./...

out:
	mkdir -p ./out


# NOTE: meta-audit は Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
darwin-amd64: out/unity-meta-check-darwin-amd64 out/unity-meta-check-junit-darwin-amd64 out/unity-meta-check-github-pr-comment-darwin-amd64 out/unity-meta-autofix-darwin-amd64 out/gh-action-yaml-gen-darwin-amd64

out/unity-meta-check-darwin-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o "$@"

out/unity-meta-check-junit-darwin-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o "$@" ./tool/unity-meta-check-junit 

out/unity-meta-check-github-pr-comment-darwin-amd64: out
	GOARCH=amd64 GOOS=darwin go build -a -tags netgo -installsuffix netgo -v -o "$@" ./tool/unity-meta-check-github-pr-comment

out/unity-meta-check-meta-audit-darwin-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o "$@" ./tool/unity-meta-check-meta-audit

out/unity-meta-autofix-darwin-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o "$@" ./tool/unity-meta-autofix

out/gh-action-yaml-gen-darwin-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o "$@" ./tool/gh-action/action-yaml-gen 

out/gh-action-darwin-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o "$@" ./tool/gh-action


darwin-arm64: out/unity-meta-check-darwin-arm64 out/unity-meta-check-junit-darwin-arm64 out/unity-meta-check-github-pr-comment-darwin-arm64 out/unity-meta-autofix-darwin-arm64 out/gh-action-yaml-gen-darwin-arm64

out/unity-meta-check-darwin-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o "$@"

out/unity-meta-check-junit-darwin-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o "$@" ./tool/unity-meta-check-junit 

out/unity-meta-check-github-pr-comment-darwin-arm64: out
	GOARCH=arm64 GOOS=darwin go build -a -tags netgo -installsuffix netgo -v -o "$@" ./tool/unity-meta-check-github-pr-comment

out/unity-meta-check-meta-audit-darwin-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o "$@" ./tool/unity-meta-check-meta-audit

out/unity-meta-autofix-darwin-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o "$@" ./tool/unity-meta-autofix

out/gh-action-yaml-gen-darwin-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o "$@" ./tool/gh-action/action-yaml-gen 

out/gh-action-darwin-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o "$@" ./tool/gh-action


# NOTE: meta-audit は Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
linux-amd64: out/unity-meta-check-linux-amd64 out/unity-meta-check-junit-linux-amd64 out/unity-meta-check-github-pr-comment-linux-amd64 out/unity-meta-autofix-linux-amd64 out/gh-action-yaml-gen-linux-amd64 out/gh-action-linux-amd64

out/unity-meta-check-linux-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o "$@"

out/unity-meta-check-junit-linux-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o "$@" ./tool/unity-meta-check-junit 

out/unity-meta-check-github-pr-comment-linux-amd64: out
	GOARCH=amd64 GOOS=linux go build -a -tags netgo -installsuffix netgo -v -o "$@" ./tool/unity-meta-check-github-pr-comment

out/unity-meta-check-meta-audit-linux-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o "$@" ./tool/unity-meta-check-meta-audit

out/unity-meta-autofix-linux-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o "$@" ./tool/unity-meta-autofix

out/gh-action-yaml-gen-linux-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o "$@" ./tool/gh-action/action-yaml-gen 

out/gh-action-linux-amd64: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o "$@" ./tool/gh-action


# NOTE: meta-audit は Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
linux-arm64: out/unity-meta-check-linux-arm64 out/unity-meta-check-junit-linux-arm64 out/unity-meta-check-github-pr-comment-linux-arm64 out/unity-meta-autofix-linux-arm64 out/gh-action-yaml-gen-linux-arm64 out/gh-action-linux-arm64

out/unity-meta-check-linux-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o "$@"

out/unity-meta-check-junit-linux-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o "$@" ./tool/unity-meta-check-junit 

out/unity-meta-check-github-pr-comment-linux-arm64: out
	GOARCH=arm64 GOOS=linux go build -a -tags netgo -installsuffix netgo -v -o "$@" ./tool/unity-meta-check-github-pr-comment

out/unity-meta-check-meta-audit-linux-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o "$@" ./tool/unity-meta-check-meta-audit

out/unity-meta-autofix-linux-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o "$@" ./tool/unity-meta-autofix

out/gh-action-yaml-gen-linux-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o "$@" ./tool/gh-action/action-yaml-gen 

out/gh-action-linux-arm64: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o "$@" ./tool/gh-action


# NOTE: meta-audit は Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
windows-amd64: out/unity-meta-check-windows-amd64.exe out/unity-meta-check-junit-windows-amd64.exe out/unity-meta-check-github-pr-comment-windows-amd64.exe out/unity-meta-autofix-windows-amd64.exe out/gh-action-yaml-gen-windows-amd64.exe

out/unity-meta-check-windows-amd64.exe: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o "$@"

out/unity-meta-check-junit-windows-amd64.exe: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o "$@" ./tool/unity-meta-check-junit 

out/unity-meta-check-github-pr-comment-windows-amd64.exe: out
	GOARCH=amd64 GOOS=windows go build -a -tags netgo -installsuffix netgo -v -o "$@" ./tool/unity-meta-check-github-pr-comment

out/unity-meta-check-meta-audit-windows-amd64.exe: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o "$@" ./tool/unity-meta-check-meta-audit

out/unity-meta-autofix-windows-amd64.exe: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o "$@" ./tool/unity-meta-autofix

out/gh-action-yaml-gen-windows-amd64.exe: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o "$@" ./tool/gh-action/action-yaml-gen 

out/gh-action-windows-amd64.exe: out
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o "$@" ./tool/gh-action


# NOTE: meta-audit は Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
windows-arm64: out/unity-meta-check-windows-arm64.exe out/unity-meta-check-junit-windows-arm64.exe out/unity-meta-check-github-pr-comment-windows-arm64.exe out/unity-meta-autofix-windows-arm64.exe out/gh-action-yaml-gen-windows-arm64.exe

out/unity-meta-check-windows-arm64.exe: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -v -o "$@"

out/unity-meta-check-junit-windows-arm64.exe: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -v -o "$@" ./tool/unity-meta-check-junit 

out/unity-meta-check-github-pr-comment-windows-arm64.exe: out
	GOARCH=arm64 GOOS=windows go build -a -tags netgo -installsuffix netgo -v -o "$@" ./tool/unity-meta-check-github-pr-comment

out/unity-meta-check-meta-audit-windows-arm64.exe: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -v -o "$@" ./tool/unity-meta-check-meta-audit

out/unity-meta-autofix-windows-arm64.exe: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -v -o "$@" ./tool/unity-meta-autofix

out/gh-action-yaml-gen-windows-arm64.exe: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -v -o "$@" ./tool/gh-action/action-yaml-gen 

out/gh-action-windows-arm64.exe: out
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o "$@" ./tool/gh-action


.PHONY: all test clean darwin-amd64 linux-amd64 windows-amd64 darwin-arm64 linux-arm64 windows-arm64
