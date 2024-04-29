package initialize

import (
	"fmt"

	"mxshop-go/stock_svc/global"

	"go.uber.org/zap"
)

func initLogger() {
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
