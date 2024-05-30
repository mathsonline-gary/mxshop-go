package model

import "gorm.io/gorm"

type SoftDelete struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
