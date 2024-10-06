package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	appLogger1 := logger.WithGroup("app1")
	appLogger2 := logger.WithGroup("app2")
	appLogger1.Info("info message", "name", "admin") // time=2024-10-05T14:47:40.758+08:00 level=INFO msg="info message" app1.name=admin
	appLogger2.Info("info message", "name", "danny")
}
