package config

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	logger "github.com/sirupsen/logrus"
)

const (
	key                    = "publish-manage"
	TimeStampFormat        = "2006-01-02 15:04:05"
	maximumCallerDepth int = 25
	minimumCallerDepth int = 4
)

func InitialLogger() {
	initLoggerDefault()

	// logger
	var logLevel logger.Level = logger.InfoLevel
	if err := logLevel.UnmarshalText([]byte(DefaultInstance.GetString("logger.level"))); err != nil {
		logger.Warnf("设置日志级别失败: %v, 将使用默认[info]级别", err)
	}
	logger.SetLevel(logLevel)
	logger.SetReportCaller(true)

	logger.AddHook(&callerHook{})
	setLoggerRotateHook()
}

func initLoggerDefault() {
	DefaultInstance.SetDefault("logger.level", "debug")
	DefaultInstance.SetDefault("moduleName", "main")
}

func setLoggerRotateHook() {
	loggerDir := GetString("logger.dir")
	if loggerDir == "" {
		loggerDir = "logs"
	}
	p, _ := filepath.Abs(loggerDir)

	p = path.Join(p, GetString("moduleName"))
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if os.MkdirAll(p, os.ModePerm) != nil {
			logger.Warn("创建日志文件夹失败!")
			return
		}
	}

	rotatedNum := GetInt("logger.rotate")
	if rotatedNum <= 0 {
		rotatedNum = 7
	}
	logFileRegex := "%Y-%m-%d.log"
	rotationTime := 24 * time.Hour

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logger.TraceLevel: writer(p, logFileRegex, "trace", rotatedNum, rotationTime),
		logger.DebugLevel: writer(p, logFileRegex, "debug", rotatedNum, rotationTime),
		logger.InfoLevel:  writer(p, logFileRegex, "info", rotatedNum, rotationTime),
		logger.WarnLevel:  writer(p, logFileRegex, "warn", rotatedNum, rotationTime),
		logger.ErrorLevel: writer(p, logFileRegex, "error", rotatedNum, rotationTime),
	}, &logger.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//}, &myFormatter{})

	logger.AddHook(lfHook)
	return
}

func writer(logPath, logFileRegex, level string, rotatedNum int, rotationTime time.Duration) io.Writer {
	logFullPath := path.Join(logPath, level)
	fileSuffix := logFileRegex

	logier, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		rotatelogs.WithLinkName(logFullPath+".log"),
		rotatelogs.WithRotationCount(rotatedNum),
		//rotatelogs.WithMaxAge(rotationTime*time.Duration(rotatedNum)),
		rotatelogs.WithRotationTime(rotationTime),
	)

	if err != nil {
		panic(err)
	}
	return logier
}

// ///////////////// caller //////////////
type callerHook struct{}

func (c *callerHook) Levels() []logger.Level {
	return []logger.Level{
		logger.PanicLevel,
		logger.FatalLevel,
		logger.ErrorLevel,
	}
}

func (c *callerHook) Fire(entry *logger.Entry) error {
	entry.Data["caller"] = c.caller(entry)
	return nil
}

func (c *callerHook) caller(entry *logger.Entry) string {
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	var stackInfo []string
	for f, again := frames.Next(); again; f, again = frames.Next() {
		if strings.Contains(f.Function, key) {
			stackInfo = append(stackInfo, fmt.Sprintf("%s:%d\n", f.Function, f.Line))
		}
	}
	//sort.Sort(sort.Reverse(sort.StringSlice(stackInfo)))
	return strings.Join(stackInfo, "\n")
}
