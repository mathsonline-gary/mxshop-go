package initialize

import (
	"fmt"

	"mxshop-go/user_svc/global"

	"go.uber.org/zap"
)

func Logger() {
	fmt.Println("logger initializing...")

	var logger *zap.Logger
	if global.Config.AppConfig.Env == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)

	fmt.Println("logger initialized!")
}
