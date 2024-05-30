// Package model defines the common data model fields.
package model

type IncrementID struct {
	ID int32 `gorm:"primarykey;type:int" json:"id"`
}
