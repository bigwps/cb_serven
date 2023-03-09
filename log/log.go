package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarlogger *zap.SugaredLogger
var logger *zap.Logger

func Logger_Init(output bool) {

	// 写入输出格式
	Encoder := SetEncoder()
	var Level zapcore.Level
	if output {
		Level = SetLevel2()
	} else {
		Level = SetLevel()
	}
	// 自定义输出级别

	//  自定义写日志级别
	LevelW := SetLevelW()
	// 自定义log保存
	Writelog := SetWritelog()
	// 控制台输出
	Console := SetConsole()
	// 创建一个core
	newCore := zapcore.NewTee(
		zapcore.NewCore(Encoder, Writelog, LevelW),                                       // 写入文件
		zapcore.NewCore(Console, zapcore.AddSync(colorable.NewColorableStdout()), Level), // 控制台
	)

	logger = zap.New(newCore)

	// logger 的. Sugar () 方法来获取一个 SugaredLogger。
	sugarlogger = logger.Sugar()

	/*
		 *Tips:
			Sync 是一个强制将缓存的数据立刻刷入磁盘的命令。
			所以 GoLang 标准库中的 File 就有 Sync 函数来对应这个命令。
			因此 logger.Sync()做的事情就是对所有输出目标文件执行 Sync。

	*/

	// defer logger.Sync()

}

// 写入文件自定义格式
func SetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "log",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
}

// 控制台自定义格式
func SetConsole() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "log",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,        // 默认换行符"\n"
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // json方式输出时不会生效, 可以设置为 capital (大写)或者 lowercase（小写）
			EncodeTime:     zapcore.ISO8601TimeEncoder,       // IS08601的时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder,   // 时间序列化，Duration为经过的浮点秒数
			EncodeCaller:   zapcore.ShortCallerEncoder,       // 路径编码器,可以设置 short （package/file:line）或者 full（ /full/path/to/package/file:line）。

		})
}

// 输出info以上的日志,
func SetLevel() zapcore.Level {
	return zapcore.InfoLevel
}

// 只输出waring以上的日志
func SetLevel2() zapcore.Level {
	return zapcore.WarnLevel
}

// 只保存error以上的日志,
func SetLevelW() zapcore.Level {
	return zapcore.WarnLevel
}

// 这里可以使用lumberjack库切割
func SetWritelog() zapcore.WriteSyncer {
	file, _ := os.Create("./zap.log")
	return zapcore.AddSync(file)
}

// 测试
func TestLog(error string) {
	//sugarlogger.Errorf(error)
	logger.Error(error)
}

// 错误日志
func ErrorLog(errin string) {
	sugarlogger.Errorf(fmt.Sprintf(errin))
}

// func Paniclog(out string) {
// 	sugarlogger.Panicf(fmt.Sprint(out))
// }

// 警告日志也兼职漏洞记录
func WarLog(errin string) {
	sugarlogger.Warnf(fmt.Sprintf(errin))
}

// info日志
func InfoLog(infoin string) {
	sugarlogger.Infof(fmt.Sprintf(infoin))
}

// 绿输出
func GreenOut(strin string) {
	color.Green(fmt.Sprintf(strin))
}

// 红输出
func RedOut(val string) {
	color.Red(fmt.Sprintf(val))
}

func Xc7Log2() {
	color.Red(`

____ ___      ____ ____ ____ _  _ ____ _  _ 
|    |__]     [__  |___ |__/ |  | |___ |\ | 
|___ |__] ___ ___] |___ |  \  \/  |___ | \| 										
	
	
◢◤ FFFFFFFFFFF ◢◤ version 1.0 ◢◤
`)
}

// test
// func HttpGet_log_test(url string) {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		sugarlogger.Errorf("连接出现了错误", url, err)

// 	} else {
// 		sugarlogger.Warnf("成功连接到地址", url, res.Status)
// 	}
// 	// defer后面必须是函数
// 	defer res.Body.Close()
// }
