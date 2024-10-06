package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	//log := zerolog.New(os.Stdout)
	//log.Info().Str("content", "Hello world").Int("count", 3).Msg("TestContextualLogger")
	//
	//// 添加上下文 (文件名/行号/字符串)
	//log = log.With().Caller().Str("foo", "bar").Logger()
	//log.Info().Msg("Hello wrold")

	sampled := log.Sample(&zerolog.BurstSampler{
		Burst:       5,
		Period:      1 * time.Second,
		NextSampler: &zerolog.BasicSampler{N: 20},
	})
	for i := 0; i <= 50; i++ {
		sampled.Info().Msgf("logged messages : %2d ", i)

	}
}
