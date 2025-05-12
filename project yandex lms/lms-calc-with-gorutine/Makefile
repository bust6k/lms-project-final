.PHONY: all grpc
all: grpc

grpc:
	protoc -I pkg/proto proto_file.proto --go_out=pkg/proto --go_opt=paths=source_relative --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative



