package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

var (
	_ driver.Valuer = (*StringList)(nil)
	_ sql.Scanner   = (*StringList)(nil)
)

type StringList []string

func (sl *StringList) Value() (driver.Value, error) {
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
