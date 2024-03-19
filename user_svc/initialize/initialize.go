package initialize

import (
	"mxshop-go/user_svc/global"
	"mxshop-go/user_svc/utils"
)

func Init() {
	// Initialize logger
	initLogger()

	// Initialize config
	initConfig()

	// Initialize database
	initDB()

	// Set app port
	if global.ServerConfig.AppConfig.Env != "production" {
		if port, err := utils.GetIdlePort(); err == nil {
			global.ServerConfig.AppConfig.Port = port
		}
	}
}
