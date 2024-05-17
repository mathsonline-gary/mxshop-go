package data

import (
	"fmt"

	"github.com/zycgary/mxshop-go/order_svc/config"
	"github.com/zycgary/mxshop-go/order_svc/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(conf config.DB, logger logger.Interface) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&model.CartItem{}, &model.Order{}, &model.OrderItem{}); err != nil {
		zap.S().Fatal("AutoMigrate tables failed: ", err)
	}
	return db
}
