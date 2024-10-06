package main

import (
	"context"
	"log/slog"
	"net"
	"os"
)

func main() {
	// 修改默认 logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	//slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
	//	AddSource:   true,            // 记录日志位置
	//	Level:       slog.LevelDebug, // 设置日志级别
	//	ReplaceAttr: nil,
	//})))

	slog.Info("hello", "name", "Danny")

	// 不是强类型:输出日志结果  !BADKEY="use of closed network connection"，提示我们 key/value 数量不匹配
	slog.Error("oops", net.ErrClosed, "status", 500)

	// 使用强类型:它限制只能传递 slog.String、slog.Int 这种强类型，如果传递普通字符串，则编译不通过。
	slog.LogAttrs(context.Background(), slog.LevelError, "oops",
		slog.Int("status", 500), slog.Any("err", net.ErrClosed))
}
