all: darwin-amd64 linux-amd64 windows-amd64


clean:
	go clean
	git clean -fdx ./dist


test:
	go test ./...


darwin-amd64:
	mkdir -p ./dist
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o ./dist/unity-meta-check-darwin-amd64
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o ./dist/unity-meta-check-junit-darwin-amd64 ./tool/unity-meta-check-junit 
	GOARCH=amd64 GOOS=darwin go build -a -tags netgo -installsuffix netgo -v -o ./dist/unity-meta-check-github-pr-comment-darwin-amd64 ./tool/unity-meta-check-github-pr-comment
	# NOTE: Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
	# CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o ./dist/unity-meta-check-meta-audit-darwin-amd64 ./tool/unity-meta-check-meta-audit
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o ./dist/unity-meta-autofix-darwin-amd64 ./tool/unity-meta-autofix


linux-amd64:
	mkdir -p ./dist
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o ./dist/unity-meta-check-linux-amd64
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o ./dist/unity-meta-check-junit-linux-amd64 ./tool/unity-meta-check-junit
	GOARCH=amd64 GOOS=linux go build -a -tags netgo -installsuffix netgo -v -o ./dist/unity-meta-check-github-pr-comment-linux-amd64 ./tool/unity-meta-check-github-pr-comment
	# NOTE: Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
	# CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o ./dist/unity-meta-check-meta-audit-linux-amd64 ./tool/unity-meta-check-meta-audit
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o ./dist/unity-meta-autofix-linux-amd64 ./tool/unity-meta-autofix


windows-amd64:
	mkdir -p ./dist
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o ./dist/unity-meta-check-windows-amd64.exe
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o ./dist/unity-meta-check-junit-windows-amd64.exe ./tool/unity-meta-check-junit
	GOARCH=amd64 GOOS=windows go build -a -tags netgo -installsuffix netgo -v -o ./dist/unity-meta-check-github-pr-comment-windows-amd64.exe ./tool/unity-meta-check-github-pr-comment
	# NOTE: Debug 用ツールなので必要になったら生成してください（生成されちゃうと Releases へあげるときに間引かないといけなくてめんどい）
	# CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o ./dist/unity-meta-check-meta-audit-windows-amd64.exe ./tool/unity-meta-check-meta-audit
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o ./dist/unity-meta-autofix-windows-amd64.exe ./tool/unity-meta-autofix


.PHONY: all test clean darwin-amd64 linux-amd64 windows-amd64
