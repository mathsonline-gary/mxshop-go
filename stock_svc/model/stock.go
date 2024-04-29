package model

type Stock struct {
	ID
	ProductID int32 `gorm:"type:int" json:"product_id"`
	Quantity  int32 `gorm:"type:int" json:"quantity"`
	Version   int32 `gorm:"type:int" json:"version"`
	Timestamps
	SoftDelete
}

type StockHistory struct {
	UserID    int32 `gorm:"type:int" json:"user_id"`
	OrderID   int32 `gorm:"type:int" json:"order_id"`
	ProductID int32 `gorm:"type:int" json:"product_id"`
	Quantity  int32 `gorm:"type:int" json:"quantity"`
	Status    int32 `gorm:"type:int" json:"status"`
}
