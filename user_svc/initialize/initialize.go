package initialize

import (
	"mxshop-go/user_svc/global"
	"mxshop-go/user_svc/utils"
)

func Init() {
	// Initialize config
	initConfig()

	// Initialize logger
	initLogger()

	// Initialize database
	initDB()

	// Set random app port for local development
	if global.Config.AppConfig.Env == "local" {
		if port, err := utils.GetIdlePort(); err == nil {
			global.Config.AppConfig.Port = port
		}
	}
}
