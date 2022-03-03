package logger

import (
	"errors"

	"github.com/go-void/portal/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	errInvalidMode = errors.New("invalid log mode")
)

type Logger struct {
	*zap.Logger

	closer func()
}

func New(c config.LogOptions) (*Logger, error) {
	if !c.Enabled {
		return &Logger{
			Logger: zap.NewNop(),
			closer: func() {},
		}, nil
	}

	ws, closer, err := zap.Open(c.Outputs...)
	if err != nil {
		return nil, err
	}

	level, err := zapcore.ParseLevel(c.Level)
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

func (l *Logger) Close() {
	l.closer()
}
