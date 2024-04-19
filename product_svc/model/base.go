package model

import (
	"database/sql"
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

var _ driver.Valuer = (*StringList)(nil)
var _ sql.Scanner = (*StringList)(nil)

func (sl StringList) Value() (driver.Value, error) {
	if sl == nil {
		return nil, nil
	}
	return json.Marshal(sl)
}

func (sl *StringList) Scan(value interface{}) error {
	if value == nil {
		*sl = []string{}
		return nil
	}
	return json.Unmarshal(value.([]byte), sl)
}
