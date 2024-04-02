package v1

import "context"

type UserDO struct {
	Nickname string `json:"name,omitempty"`
}

type UserDOList struct {
	Total int64     `json:"total,omitempty"`
	Data  []*UserDO `json:"data"`
}

type ListMeta struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty" form:"page_size"`
}

type UserRepository interface {
	Index(context.Context, ListMeta) (*UserDOList, error)
}
