# Makefile

protoc:
	protoc \
		-I./proto \
		-I/usr/local/include \
		-I$$(go env GOPATH)/src/github.com/googleapis/googleapis \
		--experimental_allow_proto3_optional \
		--go_out=./protogen/golang \
		--go_opt=paths=source_relative \
		--go-grpc_out=./protogen/golang \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./protogen/golang \
		--grpc-gateway_opt=paths=source_relative \
		proto/orders/*.proto \
		proto/product/*.proto \
		proto/google/api/*.proto \
		proto/google/type/*.proto \
		proto/google/protobuf/*.proto
