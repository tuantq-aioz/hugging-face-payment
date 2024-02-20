BIN_FILE_NAME=payment-host

run:
	@go run cmd/grpc/*.go
build:
	@cd cmd/grpc && go build -o ../../bin/$(BIN_FILE_NAME) .
proto:
	@protoc -I=internal/proto --go_out=internal/proto \
	--go-grpc_out=internal/proto \
	internal/proto/*.proto 
.PHONY: proto
