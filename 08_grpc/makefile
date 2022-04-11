# makefile 给 以下grpc使用

# 学习：08_grpc/01_grpc_helloworld
.PHONY: helloworld
helloworld:
	protoc --proto_path=. --go_out=. --go-grpc_out=. ./08_grpc/01_grpc_helloworld/proto/*.proto

# 学习：08_grpc/11_protoc_plugin
# gogoprotobuf有两个插件可以使用
  #1.protoc-gen-gogo：和protoc-gen-go生成的文件差不多，性能也几乎一样(稍微快一点点)
  #2.protoc-gen-gofast：生成的文件更复杂，性能也更高(快5-7倍)
#go get github.com/gogo/protobuf/protoc-gen-gofast
.PHONY: gofast
gofast:
	protoc --gofast_out=plugins=grpc:. ./11_protoc_gogofast/proto/*.proto


.PHONY: grpcgateway
grpcgateway:
	protoc  --proto_path=. --go_out=. --go-grpc_out=.  ./10_grpc_gateway/proto/google/api/*.proto
	protoc --proto_path=. --go_out=. --go-grpc_out=.  --grpc-gateway_out=.  ./10_grpc_gateway/proto/helloworld/*.proto


# 学习生成protobuf 使用：08_grpc/proto
.PHONY: dir_import_proto
dir_import_proto:
# 加@不打印shell
	protoc --proto_path=. --go_out=. --go-grpc_out=. ./16_import_proto/dir_import/*.proto
	#protoc --proto_path=. --go_out=plugins=grpc:. ./08_grpc/proto/dir_import/*.proto

# --------------------------------
# error生成：08_grpc/07_grpc_error
.PHONY: error
error:
# 使用gogoprotobuf
	protoc --proto_path=. --gofast_out=plugins=grpc:. ./08_grpc/07_grpc_error/proto/*.proto



.PHONY: plugin
plugin:
	# 第一步：message不生成grpc文件
	protoc --proto_path=. --go_out=.  --go-grpc_out=.  ./08_grpc/15_customized_protobuf_plugin/plugin_protobuf/*.proto
	protoc --proto_path=. --go_out=. --go-grpc_out=.  ./08_grpc/15_customized_protobuf_plugin/helloworld_protobuf/*.proto
	# 	需要手动修改pb文件的引入 _ "go_package_example/08_grpc/15_customized_protobuf_plugin/plugin_protobuf"





# 没有安装protoc-gen-go-errors
# protoc --proto_path=. --go_out=. ./08_grpc/07_grpc_error/proto/*.proto



# 安装protoc-gen-go-errors,使用 --go-errors_out=
# protoc --proto_path=. --go_out=. --go-errors_out=. ./08_grpc/07_grpc_error/proto/*.proto


# 使用绝对路径，从而忽略掉 proto 文件中的 go_package 路径
# protoc --proto_path=. --go_out=. --go-errors_out=paths=source_relative:. ./08_grpc/07_grpc_error/proto/*.proto

# ----------------------------------

