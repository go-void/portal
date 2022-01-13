package logger

import (
	"errors"

	"github.com/go-void/portal/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	errInvalidLevel = errors.New("invalid log level")
	errInvalidMode  = errors.New("invalid log mode")
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

type Logger struct {
	*zap.Logger

	closer func()
}

func New(c config.LogOptions) (*Logger, error) {
	ws, closer, err := zap.Open(c.Outputs...)
	if err != nil {
		return nil, err
	}

	level, err := GetLevel(c.Level)
	if err != nil {
		return nil, err
	}

	var encConf zapcore.EncoderConfig
	switch c.Mode {
	case "dev", "development":
		encConf = zap.NewDevelopmentEncoderConfig()
	case "", "prod", "production":
		encConf = zap.NewProductionEncoderConfig()
	default:
		return nil, errInvalidMode
	}

	enc := zapcore.NewJSONEncoder(encConf)
	core := zapcore.NewCore(enc, ws, level)

	return &Logger{
		Logger: zap.New(core),
		closer: closer,
	}, nil
}

func GetLevel(level string) (zapcore.Level, error) {
	if l, ok := levelMap[level]; ok {
		return l, nil
	}
	return -2, errInvalidLevel
}

func (l *Logger) Close() {
	l.closer()
}
