package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

type Wrapper struct {
	*zap.SugaredLogger
}

func New(env string) Logger {
	var baseLogger *zap.Logger
	var err error

	switch env {
	case "prod":
		baseLogger, err = zap.NewProduction()
	default:
		baseLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("Error creating logger", zap.Error(err))
	}

	return &Wrapper{baseLogger.Sugar()}
}
