package model

//查询实体
type FindByJobName struct {
	JobName string `bson:"jobName"` //任务名
}

// 创建实体
type TimePrint struct {
	StartTime int64 `bson:"startTime"` //开始时间
	EndTime   int64 `bson:"endTime"`   //结束时间
}
type LogRecord struct {
	JobName string    `bson:"jobName"` //任务名
	Command string    `bson:"command"` //shell命令
	Err     string    `bson:"err"`     //脚本错误
	Content string    `bson:"content"` //脚本输出
	Tp      TimePrint //执行时间
}

//更新实体
type UpdateByJobName struct {
	Command string `bson:"command"` //shell命令
	Content string `bson:"content"` //脚本输出
}

type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"tp.startTime"`
}

//startTime小于某时间，使用这种方式可以对想要进行的操作($set、$group等)提前定义
type TimeBeforeCond struct {
	BeforeTime int64 `bson:"$lt"`
}
