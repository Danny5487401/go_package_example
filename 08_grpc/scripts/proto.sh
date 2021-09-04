#!/usr/bin/env bash

# 生成golang代码
# 工具 gogoprotobuf
# protoc --proto_path=${GOPATH}/src:. --gogofaster_out=plugins=grpc:. 04_jsonpb/proto/member.proto

#protobuf
protoc  --gofast_out=plugins=grpc:./01_grpc_helloworld/proto ./01_grpc_helloworld/proto/helloworld.proto
protoc-go-inject-tag -input=./01_grpc_helloworld/proto/helloworld.pb.go

