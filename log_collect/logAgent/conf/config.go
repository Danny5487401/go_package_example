package conf

type AppConf struct {
	KafkaConf `ini:"kafka"`
	TailLogConf `ini:"tailLog"`  // 后期从etcd中获取
	EtcdConf `ini:"etcd"`

}
type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int `ini:"timeout"`
	Key string `ini:"log_collect_key"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	ChanSize int `ini:"chan_max_size"`
	Topic string `ini:"topic"`  // 后期从etcd中获取
}

// 	----used----
type TailLogConf struct {
	FilePath string `ini:"file_path"`  // 后期从etcd中获取
}