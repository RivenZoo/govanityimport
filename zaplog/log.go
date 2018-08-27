package zaplog

import (
	"go.uber.org/zap"
	"fmt"
)

const (
	LogEnvProd   = "prod"
	LogEnvStage  = "stage"
	LogEnvGlobal = "global"
)

var (
	logger *zap.Logger
	sugarLogger *zap.SugaredLogger
	logEnv string
)

func InitLogger(env string) error {
	if logEnv != "" {
		return fmt.Errorf("init logger twice, last env is %s\n", logEnv)
	}

	var err error
	switch env {
	case LogEnvStage:
		logger, err = zap.NewDevelopment()
	case LogEnvProd:
		logger, err = zap.NewProduction()
	default:
		logger = zap.L()
	}
	if err == nil {
		logEnv = env
		sugarLogger = logger.Sugar()
	}
	return err
}

func GetLogger() *zap.Logger {
	return logger
}

func GetSugarLogger() *zap.SugaredLogger {
	return sugarLogger
}

func Close() {
	if logger != nil {
		logger.Sync()
	}
}
