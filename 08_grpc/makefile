
# 学习生成protobuf 使用：08_grpc/proto
.PHONY: proto
proto:
# 加@不打印shell
	#protoc --proto_path=. --go_out=. --go-grpc_out=. ./proto/dir_import/*.proto
	protoc --proto_path=. --go_out=plugins=grpc:. ./proto/dir_import/*.proto

# --------------------------------
# error生成：08_grpc/07_grpc_error
.PHONY: error
error:
# 使用gogoprotobuf
	protoc --proto_path=. --gofast_out=plugins=grpc:. ./07_grpc_error/proto/*.proto

# 没有安装protoc-gen-go-errors
# protoc --proto_path=. --go_out=. ./07_grpc_error/proto/*.proto



# 安装protoc-gen-go-errors,使用 --go-errors_out=
# protoc --proto_path=. --go_out=. --go-errors_out=. ./07_grpc_error/proto/*.proto


# 使用绝对路径，从而忽略掉 proto 文件中的 go_package 路径
# protoc --proto_path=. --go_out=. --go-errors_out=paths=source_relative:. ./07_grpc_error/proto/*.proto

# ----------------------------------

