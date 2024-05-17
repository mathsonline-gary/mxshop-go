package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zycgary/mxshop-go/order_svc/global"
	"github.com/zycgary/mxshop-go/order_svc/model"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() {
	fmt.Println("DB initializing...")

	c := global.Config.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)

	logLevel := logger.Silent
	if global.Config.App.Debug {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
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
	if err := global.DB.AutoMigrate(&model.CartItem{}, &model.Order{}, &model.OrderItem{}); err != nil {
		zap.S().Fatal("AutoMigrate tables failed: ", err)
	}

	fmt.Println("DB initialized!")
}
