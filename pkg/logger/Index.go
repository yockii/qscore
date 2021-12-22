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

var defaultLogger *logger

func init() {
	defaultLogger = &logger{
		logrus.New(),
	}
	defaultLogger.Out = os.Stdout
	defaultLogger.SetLevel("debug")
	defaultLogger.AddHook(&GetCallerHook{Field: "caller"})
	//defaultLogger.SetReportCaller(true)
	defaultLogger.SetLogDir("logs", 30)
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
	p, _ := filepath.Abs(dir)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if os.MkdirAll(p, os.ModePerm) != nil {
			l.Warn("创建日志文件夹失败!")
			return
		}
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer(p, "debug", rotatedNum),
		logrus.InfoLevel:  writer(p, "info", rotatedNum),
		logrus.WarnLevel:  writer(p, "warn", rotatedNum),
		logrus.ErrorLevel: writer(p, "error", rotatedNum),
	}, &logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//}, &myFormatter{})
	l.AddHook(lfHook)
	return
}

func writer(logPath, level string, rotatedNum int) io.Writer {
	logFullPath := path.Join(logPath, level)
	fileSuffix := time.Now().Local().Format("2006-01-02") + ".log"

	logier, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		//rotatelogs.WithLinkName(logFullPath),
		rotatelogs.WithRotationCount(rotatedNum),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return logier
}

func SetLevel(level string) error {
	return defaultLogger.SetLevel(level)
}
func SetReportCaller(enable bool) {
	defaultLogger.SetReportCaller(enable)
}
func SetLogDir(dir string, rotateNum int) {
	defaultLogger.SetLogDir(dir, rotateNum)
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
	defaultLogger.Fatal(args...)
}
func Fatalln(args ...interface{}) {
	defaultLogger.Fatalln(args...)
}
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}
func Panicln(args ...interface{}) {
	defaultLogger.Panicln(args...)
}
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}
func Errorln(args ...interface{}) {
	defaultLogger.Errorln(args...)
}
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}
func Warnln(args ...interface{}) {
	defaultLogger.Warnln(args...)
}
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}
func Infoln(args ...interface{}) {
	defaultLogger.Infoln(args...)
}
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}
func Debugln(args ...interface{}) {
	defaultLogger.Debugln(args...)
}
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}
func Traceln(args ...interface{}) {
	defaultLogger.Traceln(args...)
}
func Tracef(format string, args ...interface{}) {
	defaultLogger.Tracef(format, args...)
}
