package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialization() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic("failed init logger: " + err.Error())
	}

	zap.ReplaceGlobals(logger)
}

func Close() {
	_ = zap.L().Sync()
}
