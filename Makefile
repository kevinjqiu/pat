.PHONY: test

test: generate
	go test -cover -v github.com/kevinjqiu/pat/pkg/...

build: generate
	go build -o pat

generate: schema
	cd pkg && go-bindata -pkg pkg schema

schema:
	cd pkg/schema && yaml2json < schema.yaml | jq . > schema.json
