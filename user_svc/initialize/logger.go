package initialize

import "go.uber.org/zap"

func initLogger() {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)
}
