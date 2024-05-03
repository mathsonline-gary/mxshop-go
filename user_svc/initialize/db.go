package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"mxshop-go/user_svc/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() {
	fmt.Println("DB initializing...")

	c := global.Config.DBConfig

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// Auto migrate models
	//if err := global.DB.AutoMigrate(&model.User{}); err != nil {
	//	panic(err)
	//}

	fmt.Println("DB initialized!")
}
