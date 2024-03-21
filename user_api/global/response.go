package global

import "time"

type Birthday time.Time

func (b Birthday) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(b).Format("2006-01-01") + `"`), nil
}

type UserResponse struct {
	ID       uint64   `json:"id"`
	Nickname string   `json:"nickname"`
	Birthday Birthday `json:"birthday"`
	Gender   int32    `json:"gender"`
	Mobile   string   `json:"mobile"`
}
