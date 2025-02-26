// pkg/logger/logger.go

package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger 初始化日志
func InitLogger(mode string) {
    config := zap.NewProductionConfig()
    
    // 开发模式下使用开发配置
    if mode != "release" {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    // 配置输出
    config.OutputPaths = []string{"stdout", "logs/app.log"}
    config.ErrorOutputPaths = []string{"stderr", "logs/error.log"}
    
    // 设置日志格式
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    config.EncoderConfig.StacktraceKey = "stacktrace"

    var err error
    Log, err = config.Build(zap.AddCallerSkip(1))
    if err != nil {
        panic(err)
    }
}

// Info wrapper for Log.Info
func Info(msg string, fields ...zap.Field) {
    Log.Info(msg, fields...)
}

// Error wrapper for Log.Error
func Error(msg string, fields ...zap.Field) {
    Log.Error(msg, fields...)
}

// Debug wrapper for Log.Debug
func Debug(msg string, fields ...zap.Field) {
    Log.Debug(msg, fields...)
}

// Warn wrapper for Log.Warn
func Warn(msg string, fields ...zap.Field) {
    Log.Warn(msg, fields...)
}

// Fatal wrapper for Log.Fatal
func Fatal(msg string, fields ...zap.Field) {
    Log.Fatal(msg, fields...)
}

// DPanic wrapper for Log.DPanic
func DPanic(msg string, fields ...zap.Field) {
    Log.DPanic(msg, fields...)
}

// Panic wrapper for Log.Panic
func Panic(msg string, fields ...zap.Field) {
    Log.Panic(msg, fields...)
}

// WithFields adds fields to log
func WithFields(fields ...zap.Field) *zap.Logger {
    return Log.With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
    return Log.Sync()
}

// 创建日志目录
func init() {
    err := os.MkdirAll("logs", 0755)
    if err != nil {
        panic(err)
    }
}