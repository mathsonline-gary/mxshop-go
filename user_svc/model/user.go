package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index:idx"`
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"index;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	Nickname string     `gorm:"type:varchar(20);"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   int32      `gorm:"type:tinyint;default:0"`
	Role     uint32     `gorm:"type:tinyint;default:1"`
}
