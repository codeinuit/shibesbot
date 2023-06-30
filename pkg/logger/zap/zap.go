package zap

import "go.uber.org/zap"

type ZapSugarLogger struct {
	logger *zap.SugaredLogger
}

func NewSugarLogger() ZapSugarLogger {
	logger, _ := zap.NewProduction()
	return ZapSugarLogger{
		logger: logger.Sugar(),
	}
}

func (z ZapSugarLogger) Info(v ...any) {
	z.logger.Info(v)
}

func (z ZapSugarLogger) Warn(v ...any) {
	z.logger.Info(v)
}

func (z ZapSugarLogger) Error(v ...any) {
	z.logger.Info(v)
}
