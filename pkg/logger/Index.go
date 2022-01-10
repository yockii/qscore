package logger

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type logger struct {
	*logrus.Logger
}

var DefaultLogger *logger

func init() {
	DefaultLogger = &logger{
		Logger: logrus.New(),
	}
	DefaultLogger.Out = os.Stdout
	DefaultLogger.SetLevel("debug")
	DefaultLogger.AddHook(&GetCallerHook{Field: "caller"})
	//DefaultLogger.SetReportCaller(true)
	//DefaultLogger.SetLogDir("logs", 30)
}

func (l *logger) SetLevel(level string) error {
	var loglevel logrus.Level
	if err := loglevel.UnmarshalText([]byte(level)); err != nil {
		l.Errorf("设置日志级别失败: %v", err)
		return err
	}
	l.Logger.SetLevel(loglevel)
	return nil
}

func (l *logger) SetReportCaller(enable bool) {
	l.Logger.SetReportCaller(enable)
}

func (l *logger) SetLogDir(dir string, rotatedNum int) {
	l.SetLogDirWithLogFileRegexAndRotateTime(dir, "%Y-%m-%D.log", rotatedNum, 24*time.Hour)
}
func (l *logger) SetLogDirWithLogFileRegex(dir, logFileRegex string, rotatedNum int) {
	l.SetLogDirWithLogFileRegexAndRotateTime(dir, logFileRegex, rotatedNum, 24*time.Hour)
}
func (l *logger) SetLogDirRotateTime(dir string, rotatedNum int, rotationTime time.Duration) {
	l.SetLogDirWithLogFileRegexAndRotateTime(dir, "%Y-%m-%D.log", rotatedNum, rotationTime)
}

func (l *logger) SetLogDirWithLogFileRegexAndRotateTime(dir, logFileRegex string, rotatedNum int, rotationTime time.Duration) {
	p, _ := filepath.Abs(dir)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if os.MkdirAll(p, os.ModePerm) != nil {
			l.Warn("创建日志文件夹失败!")
			return
		}
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer(p, logFileRegex, "debug", rotatedNum, rotationTime),
		logrus.InfoLevel:  writer(p, logFileRegex, "info", rotatedNum, rotationTime),
		logrus.WarnLevel:  writer(p, logFileRegex, "warn", rotatedNum, rotationTime),
		logrus.ErrorLevel: writer(p, logFileRegex, "error", rotatedNum, rotationTime),
	}, &logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//}, &myFormatter{})

	l.AddHook(lfHook)
	return
}

func writer(logPath, logFileRegex, level string, rotatedNum int, rotationTime time.Duration) io.Writer {
	logFullPath := path.Join(logPath, level)
	fileSuffix := logFileRegex

	logier, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		rotatelogs.WithLinkName(logFullPath),
		rotatelogs.WithRotationCount(rotatedNum),
		//rotatelogs.WithMaxAge(rotationTime*time.Duration(rotatedNum)),
		rotatelogs.WithRotationTime(rotationTime),
	)

	if err != nil {
		panic(err)
	}
	return logier
}

func SetLevel(level string) error {
	return DefaultLogger.SetLevel(level)
}
func SetReportCaller(enable bool) {
	DefaultLogger.SetReportCaller(enable)
}
func SetLogDir(dir string, rotateNum int) {
	DefaultLogger.SetLogDir(dir, rotateNum)
}

//////////////////////////////////////////////////////////////////////////
func (l *logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}
func (l *logger) Fatalln(args ...interface{}) {
	l.Logger.Fatalln(args...)
}
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.Logger.Panic(args...)
}
func (l *logger) Panicln(args ...interface{}) {
	l.Logger.Panicln(args...)
}
func (l *logger) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}
func (l *logger) Errorln(args ...interface{}) {
	l.Logger.Errorln(args...)
}
func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}
func (l *logger) Warnln(args ...interface{}) {
	l.Logger.Warnln(args...)
}
func (l *logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}
func (l *logger) Infoln(args ...interface{}) {
	l.Logger.Infoln(args...)
}
func (l *logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}
func (l *logger) Debugln(args ...interface{}) {
	l.Logger.Debugln(args...)
}
func (l *logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

func (l *logger) Trace(args ...interface{}) {
	l.Logger.Trace(args...)
}
func (l *logger) Traceln(args ...interface{}) {
	l.Logger.Traceln(args...)
}
func (l *logger) Tracef(format string, args ...interface{}) {
	l.Logger.Tracef(format, args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}
func Fatalln(args ...interface{}) {
	DefaultLogger.Fatalln(args...)
}
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	DefaultLogger.Panic(args...)
}
func Panicln(args ...interface{}) {
	DefaultLogger.Panicln(args...)
}
func Panicf(format string, args ...interface{}) {
	DefaultLogger.Panicf(format, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}
func Errorln(args ...interface{}) {
	DefaultLogger.Errorln(args...)
}
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}
func Warnln(args ...interface{}) {
	DefaultLogger.Warnln(args...)
}
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}
func Infoln(args ...interface{}) {
	DefaultLogger.Infoln(args...)
}
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}
func Debugln(args ...interface{}) {
	DefaultLogger.Debugln(args...)
}
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Trace(args ...interface{}) {
	DefaultLogger.Trace(args...)
}
func Traceln(args ...interface{}) {
	DefaultLogger.Traceln(args...)
}
func Tracef(format string, args ...interface{}) {
	DefaultLogger.Tracef(format, args...)
}
