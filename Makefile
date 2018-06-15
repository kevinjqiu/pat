.PHONY: test

test:
	go test -cover -v github.com/kevinjqiu/pat/pkg/...

build:
	go build -o pat
