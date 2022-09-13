package main

import (
	"fmt"
	"strconv"

	_ "go_package_example/08_grpc/15_customized_protobuf_plugin/helloworld_protobuf"
	options "go_package_example/08_grpc/15_customized_protobuf_plugin/plugin_protobuf"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func main() {
	// 遍历文件
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {

		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			// 获取服务
			service := services.Get(i)
			if serviceHandler, _ := proto.GetExtension(service.Options(), options.E_ServiceHandler).(*options.ServiceHandler); serviceHandler != nil {
				fmt.Println("--- service ---")
				fmt.Println("service name: " + string(service.FullName()))

				if serviceHandler.Authorization != nil && *serviceHandler.Authorization != "" {
					fmt.Println("use interceptor authorization: " + *serviceHandler.Authorization)
				}
				fmt.Println("--- service ---")
			}

			// 获取方法
			methods := service.Methods()
			for k := 0; k < methods.Len(); k++ {
				method := methods.Get(k)
				if methodHandler, _ := proto.GetExtension(method.Options(), options.E_MethodHandler).(*options.MethodHandler); methodHandler != nil {
					fmt.Println("--- method ---")
					fmt.Println("method name: " + string(method.FullName()))
					if methodHandler.Whitelist != nil && *methodHandler.Whitelist != "" {
						fmt.Println("use interceptor whitelist: " + *methodHandler.Whitelist)
					}

					if methodHandler.Logger != nil {
						fmt.Println("use interceptor logger: " + strconv.FormatBool(*methodHandler.Logger))
					}

					fmt.Println("--- method ---")
				}
			}
		}

		return true
	})
}
