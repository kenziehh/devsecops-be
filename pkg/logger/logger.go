package logger

import (
    "context"
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type Fields map[string]interface{}

type Logger interface {
    Debug(ctx context.Context, message string, fields ...Fields)
    Info(ctx context.Context, message string, fields ...Fields)
    Warn(ctx context.Context, message string, fields ...Fields)
    Error(ctx context.Context, message string, err error, fields ...Fields)
    Fatal(ctx context.Context, message string, err error, fields ...Fields)
}

type logger struct {
    zap *zap.Logger
}

func NewLogger() Logger {
    config := zap.NewProductionConfig()
    
    // Set log level based on environment
    if os.Getenv("ENV") == "development" {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    config.OutputPaths = []string{"stdout"}
    config.ErrorOutputPaths = []string{"stderr"}
    
    zapLogger, err := config.Build(zap.AddCallerSkip(1))
    if err != nil {
        panic(err)
    }

    return &logger{
        zap: zapLogger,
    }
}

func (l *logger) Debug(ctx context.Context, message string, fields ...Fields) {
    l.zap.Debug(message, l.fieldsToZap(fields...)...)
}

func (l *logger) Info(ctx context.Context, message string, fields ...Fields) {
    l.zap.Info(message, l.fieldsToZap(fields...)...)
}

func (l *logger) Warn(ctx context.Context, message string, fields ...Fields) {
    l.zap.Warn(message, l.fieldsToZap(fields...)...)
}

func (l *logger) Error(ctx context.Context, message string, err error, fields ...Fields) {
    zapFields := l.fieldsToZap(fields...)
    if err != nil {
        zapFields = append(zapFields, zap.Error(err))
    }
    l.zap.Error(message, zapFields...)
}

func (l *logger) Fatal(ctx context.Context, message string, err error, fields ...Fields) {
    zapFields := l.fieldsToZap(fields...)
    if err != nil {
        zapFields = append(zapFields, zap.Error(err))
    }
    l.zap.Fatal(message, zapFields...)
}

func (l *logger) fieldsToZap(fields ...Fields) []zap.Field {
    var zapFields []zap.Field
    
    for _, fieldMap := range fields {
        for key, value := range fieldMap {
            zapFields = append(zapFields, zap.Any(key, value))
        }
    }
    
    return zapFields
}