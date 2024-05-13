clean:
	go clean
	git clean -fdx ./out


test:
	go test ./...


.PHONY: test clean
