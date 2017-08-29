package clog

import (
	"github.com/plimble/zap-sentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Dev       bool   `env:"ZAP_DEV"`
	Level     string `env:"ZAP_LEVEL" validate:"zap_level"`
	DSN       string
	ModuleKey string
}

type Logger struct {
	*zap.Logger
	config Config
}

func NewLogger(config Config) (Logger, error) {
	if config.ModuleKey == "" {
		config.ModuleKey = "module"
	}

	var zcfg zap.Config
	if config.Dev {
		zcfg = zap.NewDevelopmentConfig()
	} else {
		zcfg = zap.NewProductionConfig()
	}

	if config.Level != "" {
		level := new(zapcore.Level)
		if err := level.Set(config.Level); err != nil {
			return Logger{}, err
		}
		zcfg.Level.SetLevel(*level)
	}

	var l *zap.Logger
	if config.Dev {
		var err error
		l, err = zcfg.Build()
		if err != nil {
			return Logger{}, err
		}
	} else {
		scfg := zapsentry.Configuration{DSN: config.DSN}
		sentryCore, err := scfg.Build()
		if err != nil {
			return Logger{}, err
		}
		sentryCoreFn := func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, sentryCore)
		}
		l, err = zcfg.Build(zap.WrapCore(sentryCoreFn))
		if err != nil {
			return Logger{}, err
		}
	}
	return Logger{Logger: l, config: config}, nil
}

func (l *Logger) Module(name string, fields ...zapcore.Field) *zap.Logger {
	fields = append(fields, zap.String(l.config.ModuleKey, name))
	return l.With(fields...)
}
