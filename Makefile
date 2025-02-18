# Makefile

protoc:
	cd proto && protoc --plugin=protoc-gen-go=/Users/dan/go/bin/protoc-gen-go \
		--go_out=../protogen/golang --go_opt=paths=source_relative \
		./**/*.proto
