package data

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/config"
	"github.com/zycgary/mxshop-go/pkg/common/model"
	"github.com/zycgary/mxshop-go/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type User struct {
	model.IncrementID
	model.Timestamps
	model.SoftDelete

	Email    string     `gorm:"index;unique;type:varchar(255);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	Nickname string     `gorm:"type:varchar(20);"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   int32      `gorm:"type:tinyint;default:0"`
	Role     uint32     `gorm:"type:tinyint;default:1"`
}

// NewDB creates a new GORM database connection and migrates the models.
func NewDB(conf config.DB, logger log.Logger) (*gorm.DB, error) {
	ls := log.NewSugar(logger)

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
	if err := db.AutoMigrate(&User{}); err != nil {
		ls.Fatalf("AutoMigrate: %s", err)
	}

	return db, nil
}
