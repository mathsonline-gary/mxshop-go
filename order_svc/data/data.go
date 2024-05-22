package data

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/zycgary/mxshop-go/order_svc/config"
	"github.com/zycgary/mxshop-go/order_svc/model"
	"github.com/zycgary/mxshop-go/pkg/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewGormDB(conf config.DB, logger log.Logger) (*gorm.DB, error) {
	// Init DB logger
	logLevel := gormlogger.Silent
	if logger.Level() == log.LevelDebug {
		logLevel = gormlogger.Info
	}
	dbLogger := gormlogger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	if err := db.AutoMigrate(&model.CartItem{}, &model.Order{}, &model.OrderItem{}); err != nil {
		zap.S().Fatal("AutoMigrate tables failed: ", err)
	}
	return db, nil
}
