package main

import (
	"io"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	InitLog("02_log/02_zap/04_customized_log/log/info.log", "02_log/02_zap/04_customized_log/log/error.log", zap.InfoLevel)
	defer logger.Sync()

	sugarLogger.Infof("sugarLogger name:%s", "修华师1")
	sugarLogger.Infow("sugarLogger", zap.String("name", "修华师2"))
	sugarLogger.Errorf("sugarLogger name:%s", "修华师3")
	sugarLogger.Debugf("sugarLogger name:%s", "修华师4")
	sugarLogger.Warnf("sugarLogger name:%s", "修华师5")

	logger.Info("logger", zap.String("name", "修华师6"))
	logger.Error("logger", zap.String("name", "修华师7"))
	logger.Debug("logger", zap.String("name", "修华师8"))
}

// 只能输出结构化日志，但是性能要高于 SugaredLogger
var logger *zap.Logger

// 可以输出 结构化日志、非结构化日志。性能差于 zap.Logger，
var sugarLogger *zap.SugaredLogger

// 初始化日志 logger
func InitLog(logPath, errPath string, logLevel zapcore.Level) {
	config := zapcore.EncoderConfig{
		MessageKey:   "msg",                       //结构化（json）输出：msg的key
		LevelKey:     "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/sql_squirrel_test.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	//自定义日志级别：自定义Warn级别
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	// 获取io.Writer的实现
	infoWriter := getWriter1(logPath)
	warnWriter := getWriter2(errPath)

	// 实现多个输出
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), infoLevel),                         //将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(warnWriter), warnLevel),                         //warn及以上写入errPath
		zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	sugarLogger = logger.Sugar()
}

func getWriter1(filename string) io.Writer {
	// github.com/lestrrat-go/file-rotatelogs 已经2021年 archived(不推荐)
	// 生成 rotatelogs 的Logger 实际生成的文件名 filename.YYmmddHH
	// filename是指向最新日志的链接
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24), // 切割频率 24小时
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}

func getWriter2(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,    //文件大小限制,单位MB，超过则切割
		MaxBackups: 5,     //最大文件保留数，超过就删除最老的日志文件
		MaxAge:     30,    //日志文件保留天数
		Compress:   false, //是否压缩
	}
}
