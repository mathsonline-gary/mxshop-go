package model

import "time"

type CartItem struct {
	IncrementID
	UserID    int32 `gorm:"type:int;index" json:"user_id"`
	ProductID int32 `gorm:"type:int;index" json:"product_id"`
	Quantity  int32 `gorm:"type:int" json:"quantity"`
	Selected  bool  `gorm:"type:bool" json:"selected"`
	Timestamps
	SoftDelete
}

type Order struct {
	IncrementID
	UserID          int32     `gorm:"type:int;index" json:"user_id"`
	SerialNumber    string    `gorm:"type:varchar(100);index" json:"serial_number"`
	PaymentMethod   string    `gorm:"type:varchar(20)" json:"payment_method"`
	Status          string    `gorm:"type:varchar(20)" json:"status"`
	TradingNumber   string    `gorm:"type:varchar(100)" json:"trading_number"`
	Amount          float32   `gorm:"type:float" json:"amount"`
	PaidAt          time.Time `gorm:"type:datetime" json:"paid_at"`
	ShippingAddress string    `gorm:"type:varchar(255)" json:"shipping_address"`
	ShippingName    string    `gorm:"type:varchar(100)" json:"shipping_name"`
	ShippingMobile  string    `gorm:"type:varchar(20)" json:"shipping_mobile"`
	Note            string    `gorm:"type:varchar(255)" json:"note"`
	Timestamps
	SoftDelete
}

type OrderItem struct {
	IncrementID
	OrderID           int32   `gorm:"type:int;index" json:"order_id"`
	ProductID         int32   `gorm:"type:int;index" json:"product_id"`
	ProductName       string  `gorm:"type:varchar(100)" json:"product_name"`
	ProductImage      string  `gorm:"type:varchar(255)" json:"product_image"`
	ProductUnitPrice  float32 `gorm:"type:float" json:"product_unit_price"`
	ProductTotalPrice float32 `gorm:"type:float" json:"product_price"`
	Quantity          int32   `gorm:"type:int" json:"quantity"`
	Timestamps
	SoftDelete
}
