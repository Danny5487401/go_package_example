# makefile 给 以下 grpc使用

# 1：08_grpc/01_grpc_helloworld
.PHONY: helloworld
helloworld:
	protoc --proto_path=. --go_out=. --go-grpc_out=. 01_grpc_helloworld/proto/*.proto


# 3: 08_grpc/03_stream_grpc
.PHONY: stream-grpc
stream-grpc:
	protoc --go_out=. --go-grpc_out=. 03_stream_grpc/proto/stream.proto


# 4: 08_grpc/04_jsonpb
.PHONY: jsonpb
jsonpb:
	protoc --go_out=. --go-grpc_out=. 04_jsonpb/proto/member.proto



# 7：08_grpc/07_grpc_error
.PHONY: grpc_error
grpc_error:
	# 使用 gogoprotobuf
	protoc --gofast_out=. --go-grpc_out=. 07_grpc_error/proto/*.proto

# 9: 08_grpc/09_grpc_validate/
.PHONY: grpc_validate
grpc_validate:
	# 使用 gogoprotobuf
	protoc   -I . -I 09_grpc_validate/proto --gofast_out=. --go-grpc_out=. 09_grpc_validate/proto/helloworld.proto


# 10: 08_grpc/10_grpc_gateway
.PHONY: grpcgateway
grpcgateway:
	protoc  --proto_path=. --go_out=. --go-grpc_out=.  ./10_grpc_gateway/proto_without_buf/google/api/*.proto
	protoc --proto_path=. --go_out=. --go-grpc_out=.  --grpc-gateway_out=.  ./10_grpc_gateway/proto_without_buf/helloworld/*.proto

# 15: 08_grpc/15_customized_protobuf_plugin
.PHONY: plugin
plugin:
	# 第一步：message不生成grpc文件
	protoc --proto_path=. --go_out=.  --go-grpc_out=.  ./08_grpc/15_customized_protobuf_plugin/plugin_protobuf/*.proto
	protoc --proto_path=. --go_out=. --go-grpc_out=.  ./08_grpc/15_customized_protobuf_plugin/helloworld_protobuf/*.proto
	# 	需要手动修改pb文件的引入 _ "go_package_example/08_grpc/15_customized_protobuf_plugin/plugin_protobuf"




# 16：08_grpc/16_import_proto
.PHONY: dir_import_proto
dir_import_proto:
# 加@不打印shell
	protoc --proto_path=. --go_out=. --go-grpc_out=. 16_import_proto/proto/component.proto 16_import_proto/proto/computer.proto
	#protoc --proto_path=. --go_out=plugins=grpc:. ./08_grpc/proto/dir_import/*.proto








