#!/usr/bin/env bash

# 生成golang代码
# 工具 gogoprotobuf
# protoc --proto_path=${GOPATH}/src:. --gogofaster_out=plugins=grpc:. 04_jsonpb/proto/member.proto

#protobuf
protoc  --gofast_out=plugins=grpc:./04_jsonpb/proto ./04_jsonpb/proto/member.proto
protoc-go-inject-tag -input=./04_jsonpb/proto/member.pb.go

