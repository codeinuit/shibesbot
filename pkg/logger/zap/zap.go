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

func (z ZapSugarLogger) Infof(format string, v ...any) {
	z.logger.Infof(format, v...)
}

func (z ZapSugarLogger) Warn(v ...any) {
	z.logger.Warn(v)
}

func (z ZapSugarLogger) Warnf(format string, v ...any) {
	z.logger.Warnf(format, v...)
}

func (z ZapSugarLogger) Error(v ...any) {
	z.logger.Error(v)
}

func (z ZapSugarLogger) Errorf(format string, v ...any) {
	z.logger.Errorf(format, v...)
}
