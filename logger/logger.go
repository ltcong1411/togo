package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

// GetLogger ...
func GetLogger(module string) *zap.SugaredLogger {
	logLevel := os.Getenv("LOG_LEVEL")
	runEnv := os.Getenv("RUN_ENV")
	var config zap.Config

	if strings.ToUpper(runEnv) == "PROD" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level.UnmarshalText([]byte(logLevel))
	log, _ := config.Build()

	return log.Named(module).Sugar()
}
