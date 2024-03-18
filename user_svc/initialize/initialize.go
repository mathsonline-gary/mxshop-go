package initialize

func Init() {
	// Initialize logger
	initLogger()

	// Initialize config
	initConfig()

	// Initialize database
	initDB()
}
