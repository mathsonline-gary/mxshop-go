package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type IncrementID struct {
	ID int32 `gorm:"primarykey;type:int" json:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SoftDelete struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
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
