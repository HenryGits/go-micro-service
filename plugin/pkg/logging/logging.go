/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 全局日志处理模块
**/

package logging

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	kitZapLog "github.com/go-kit/kit/log/zap"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"path/filepath"
	"time"
)

const (
	LoggerRequestId = "trace-id"
	TraceId         = "trace-id"
)

/**
 * logger: 日志对象
 * levelOut: 日志等级
 */
func SetLogging(logger *zap.Logger, logPath string, level zapcore.Level) kitlog.Logger {
	var log kitlog.Logger
	if logPath != "" {
		// default log
		log = defaultLogger(logPath)
		log = kitlog.WithPrefix(log, "ts", kitlog.TimestampFormat(func() time.Time {
			return time.Now()
		}, "2006-01-02 15:04:05"))
	} else {
		//logger = kitlog.NewLogfmtLogger(kitlog.StdlibWriter{})
		//log = term.NewLogger(os.Stdout, kitlog.NewLogfmtLogger, colorFunc())
		log = kitZapLog.NewZapSugarLogger(logger, level)
		log = kitlog.WithPrefix(log, "ts", kitlog.TimestampFormat(func() time.Time {
			return time.Now()
		}, "2006-01-02 15:04:05"))
	}
	//logger = level.NewFilter(logger, logLevel(levelOut))
	//logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)

	return log
}

func logLevel(logLevel string) (opt level.Option) {
	switch logLevel {
	case "warn":
		opt = level.AllowWarn()
	case "error":
		opt = level.AllowError()
	case "debug":
		opt = level.AllowDebug()
	case "info":
		opt = level.AllowInfo()
	case "all":
		opt = level.AllowAll()
	default:
		opt = level.AllowNone()
	}
	return opt
}

func defaultLogger(filePath string) kitlog.Logger {
	linkFile, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	writer, err := rotatelogs.New(
		linkFile+"-%Y-%m-%d",
		rotatelogs.WithLinkName(linkFile),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*365),   // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return kitlog.NewLogfmtLogger(writer)
}

func colorFunc() func(keyVals ...interface{}) term.FgBgColor {
	return func(keyVals ...interface{}) term.FgBgColor {
		for i := 0; i < len(keyVals)-1; i += 2 {
			if keyVals[i] != "level" {
				continue
			}
			val := fmt.Sprintf("%v", keyVals[i+1])
			switch val {
			case "debug":
				return term.FgBgColor{Fg: term.DarkGray}
			case "info":
				return term.FgBgColor{Fg: term.Blue}
			case "warn":
				return term.FgBgColor{Fg: term.Yellow}
			case "error":
				return term.FgBgColor{Fg: term.Red}
			case "crit":
				return term.FgBgColor{Fg: term.Gray, Bg: term.DarkRed}
			default:
				return term.FgBgColor{}
			}
		}
		return term.FgBgColor{}
	}
}
