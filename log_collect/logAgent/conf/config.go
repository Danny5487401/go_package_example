package conf

type AppConf struct {
	KafkaConf `ini:"kafka"`
	TailLogConf `ini:"tailLog"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic string `ini:"topic"`
}

type TailLogConf struct {
	FilePath string `ini:"file_path"`
}