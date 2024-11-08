package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLog struct {
	logger *zap.SugaredLogger
}

func NewZapLog(level string) (*ZapLog, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLog{
		logger: logger.Sugar(),
	}, nil
}

func (l *ZapLog) Close() {
	_ = l.logger.Sync()
}

func (l *ZapLog) WithField(key string, value any) Logger {
	return &ZapLog{
		logger: l.logger.With(zap.Any(key, value)),
	}
}

func (l *ZapLog) WithError(err error) Logger {
	return &ZapLog{
		logger: l.logger.With(zap.Any("error", err.Error())),
	}
}

func (l *ZapLog) Error(msg string) {
	l.logger.Error(msg)
}

func (l *ZapLog) Errorf(template string, args ...any) {
	l.logger.Errorf(template, args...)
}

func (l *ZapLog) Warning(msg string) {
	l.logger.Warn(msg)
}

func (l *ZapLog) Warningf(template string, args ...any) {
	l.logger.Warnf(template, args...)
}

func (l *ZapLog) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLog) Infof(template string, args ...any) {
	l.logger.Infof(template, args...)
}

func (l *ZapLog) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *ZapLog) Debugf(template string, args ...any) {
	l.logger.Debugf(template, args...)
}
