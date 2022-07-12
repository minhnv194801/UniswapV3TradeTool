package logger

import (
	"log"
	"net/url"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	logger      *zap.Logger
	serviceName string
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (l lumberjackSink) Sync() error {
	return nil
}

func NewZapLogger(initial, thereafter int, serviceName string) (*ZapLogger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Sampling = &zap.SamplingConfig{
		Initial:    initial,
		Thereafter: thereafter,
		Hook: func(e zapcore.Entry, d zapcore.SamplingDecision) {
			if d == zapcore.LogDropped {
				log.Println("log dropped, log per second exceeded")
			}
		},
	}
	err := zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumberjack.Logger{
				Filename:   u.Opaque,
				MaxSize:    128,
				MaxAge:     28,
				MaxBackups: 3,
				LocalTime:  false,
				Compress:   false,
			},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	loggerConfig.OutputPaths = append(loggerConfig.OutputPaths, "lumberjack:bot.log")
	logger, err := loggerConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: logger, serviceName: serviceName}, nil
}

func (logger *ZapLogger) Info(funcName string, msg string, attempt uint16) {
	logger.logger.Info(msg,
		zap.Int64("timestamp", time.Now().Unix()),
		zap.String("service", logger.serviceName),
		zap.String("func", funcName),
		zap.Uint16("attempt", attempt))
}

func (logger *ZapLogger) Warn(funcName string, msg string, attempt uint16) {
	logger.logger.Warn(msg,
		zap.Int64("timestamp", time.Now().Unix()),
		zap.String("service", logger.serviceName),
		zap.String("func", funcName),
		zap.Uint16("attempt", attempt))
}

func (logger *ZapLogger) Error(funcName string, msg string, attempt uint16) {
	logger.logger.Error(msg,
		zap.Int64("timestamp", time.Now().Unix()),
		zap.String("service", logger.serviceName),
		zap.String("func", funcName),
		zap.Uint16("attempt", attempt))
}
