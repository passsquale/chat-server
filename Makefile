LOCAL_BIN:=$(CURDIR)/bin

install-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	mkdir -p bin
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc

generate:
	make generate-chat-api
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png '

generate-chat-api:
	mkdir -p pkg/chat_v1
	mkdir -p pkg/swagger
	protoc --proto_path api/chat_v1 --proto_path vendor.protogen \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/chat_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--grpc-gateway_out=pkg/chat_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/chat_v1/chat.proto		

go-build:
	go build -o ./build/auth_server cmd/server/main.go

docker-build:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/alexzabolotskikh/chat_server .

up-local:
	docker-compose up -d

run-app-local: up-local
	go run cmd/server/main.go --config-path=local.env

test:
	go clean -testcache
	go test ./... -coverpkg=./internal...

test-cover:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -coverpkg=./internal...
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out
	go tool cover -func=./coverage.out | grep "total"

vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p  vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
	fi