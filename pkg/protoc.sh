#!/bin/bash
#
# Make sure protocol compiler plugins for Go are installed
#
# $ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#

protoc --go_out job --go_opt paths=source_relative --go-grpc_out job --go-grpc_opt paths=source_relative job.proto
