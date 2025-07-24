package x

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Global Logger use Zap log
var Xlog *zap.Logger

// Config Configuration of log
// ErrorFilePath: error log file path
// InfoFilePath: info log file path
// MaxSize: megabytes
// MaxBackups: number of log files
// MaxAge:days
const (
	LogPath       = "./log" //"N:/PLATFORM/SERVER/Config/Log"
	ErrorFileName = "x-err.log"
	InfoFileName  = "x-info.log"
	MaxSize       = 50
	MaxBackups    = 10
	MaxAge        = 365
)

// TODO : #  1. 设默认值 2. 加载默认配置(会引入viper依赖?)去覆盖

func init() {
	fmt.Println("init xlog Start")

	if err := InitXLogger(); err != nil {
		panic(err)
	}

	{
		Xlog.Error("create log file error success")
		Xlog.Info("create log file info success")
	}

	fmt.Println("init xlog End")
}

func InitXLogger() error {
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderJSON := zapcore.NewJSONEncoder(encoder)

	errWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   LogPath + "/" + ErrorFileName,
		MaxSize:    MaxSize, // megabytes
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge, // days
	})
	infoWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   LogPath + "/" + InfoFileName,
		MaxSize:    MaxSize, // megabytes
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge, // days
	})

	errPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	errCore := zapcore.NewCore(encoderJSON, errWriteSyncer, errPriority)
	infoCore := zapcore.NewCore(encoderJSON, infoWriteSyncer, infoPriority)
	Xlog = zap.New(zapcore.NewTee(errCore, infoCore), zap.AddCaller())

	return nil
}

// Field is an alias for zap.Field. Aliasing this type dramatically
type Field = zap.Field

// Int constructs a field with the given key and value.
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return zap.String(key, val)
}

// // Debug logs a message at DebugLevel
// func Debug(msg string, fields ...Field) {
// 	Xlog.Debug(msg, fields...)
// }

// // Info logs a message at InfoLevel
// func Info(msg string, fields ...Field) {
// 	Xlog.Info(msg, fields...)
// }

// // Warn logs a message at WarnLevel
// func Warn(msg string, fields ...Field) {
// 	Xlog.Warn(msg, fields...)
// }

// // Error logs a message at ErrorLevel
// func Error(msg string, fields ...Field) {
// 	Xlog.Error(msg, fields...)
// }

// // ErrErr logs a message at ErrorLevel
// func ErrErr(msg string, err error) {
// 	Xlog.Error(msg, String("error", err.Error()))
// }

// // DPanic logs a message at DPanicLevel
// func DPanic(msg string, fields ...Field) {
// 	Xlog.DPanic(msg, fields...)
// }

// // Panic logs a message at PanicLevel
// func Panic(msg string, fields ...Field) {
// 	Xlog.Panic(msg, fields...)
// }

// // Fatal logs a message at FatalLevel
// func Fatal(msg string, fields ...Field) {
// 	Xlog.Fatal(msg, fields...)
// }
