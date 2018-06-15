.PHONY: test

test:
	go test -v github.com/kevinjqiu/pat/pkg/...

build:
	go build -o pat
