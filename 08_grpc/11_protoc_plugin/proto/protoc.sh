# gogoprotobuf有两个插件可以使用

  #1.protoc-gen-gogo：和protoc-gen-go生成的文件差不多，性能也几乎一样(稍微快一点点)
  #2.protoc-gen-gofast：生成的文件更复杂，性能也更高(快5-7倍)

#go get github.com/gogo/protobuf/protoc-gen-gofast
protoc --gofast_out=plugins=grpc:. helloworld.proto