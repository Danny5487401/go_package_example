package conf

type AppConf struct {
	KafkaConf `ini:"kafka"`
	TailLogConf `ini:"tailLog"`  // 后期从etcd中获取
	EtcdConf `ini:"etcd"`

}
type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int `ini:"timeout"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic string `ini:"topic"`  // 后期从etcd中获取
}

// 	----used----
type TailLogConf struct {
	FilePath string `ini:"file_path"`  // 后期从etcd中获取
}