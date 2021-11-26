package hooks

import (
	"context"
	"github.com/go-redis/redis/v7"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "redis: ", log.LstdFlags|log.Lshortfile)

func SetLogger(logLogger *log.Logger) {
	logger = logLogger
}

func NewLogHook() redis.Hook {
	return &logHook{}
}

type logHook struct{}

func (h *logHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *logHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	logger.Printf("cmd: %v \n", cmd.String())
	return nil
}

func (h *logHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *logHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	args := make([][]string, 0, len(cmds))
	for _, cmd := range cmds {
		args = append(args, []string{cmd.String()})
	}
	logger.Printf("pipline: %v \n", args)
	return nil
}
