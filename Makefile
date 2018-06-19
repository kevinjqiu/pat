.PHONY: test

test: generate
	go test -cover -v github.com/kevinjqiu/pat/pkg/...

build: generate
	go build -o pat

generate: bindata schema
	cd pkg && go-bindata -pkg pkg schema

schema: yaml2json
	cd pkg/schema && yaml2json < schema.yaml | jq . > schema.json

bindata:
	go get -u github.com/jteeuwen/go-bindata

yaml2json:
	go get -u github.com/bronze1man/yaml2json
