package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"log"
)

// 组件定义
type Config struct {
	// 配置信息可以放在这里
	Name string
	Age  int
}

type Service struct {
	config *Config
}

func NewConfig() *Config {
	// 配置信息的初始化逻辑
	return &Config{
		Name: "Danny",
		Age:  20,
	}

}

func NewService(config *Config) *Service {
	return &Service{config: config}
}

// 应用启动时执行的操作
func (s *Service) OnStart() {
	fmt.Println("Service started with config:", s.config)
}

// 应用停止时执行的操作
func (s *Service) OnStop() {
	fmt.Println("Service stopped")
}

func main() {
	app := fx.New(
		// 提供 Config 的构建方法
		fx.Provide(NewConfig),
		// 提供 Service 的构建方法，并注入 Config 作为依赖
		fx.Provide(NewService),
		// 注册 Service 的 OnStart 和 OnStop 生命周期事件
		fx.Invoke(func(lifecycle fx.Lifecycle, service *Service) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					service.OnStart()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					service.OnStop()
					return nil
				},
			})

		}),
	)

	ctx := context.Background()
	// 启动应用
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer app.Stop(ctx)

	// 在这里可以进行更多的应用逻辑操作
	fmt.Println("Application is running...")
}
