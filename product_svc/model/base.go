package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32 `gorm:"primarykey;type:int"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type StringList []string

func (sl *StringList) Value() (driver.Value, error) {
	return json.Marshal(sl)
}

func (sl *StringList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), sl)
}
