package model

type Category struct {
	BaseModel
	Name                 string      `gorm:"type:varchar(20);not null" json:"name"`
	Level                int32       `gorm:"type:int;not null;default:1" json:"level"`
	VisibleInTab         bool        `gorm:"not null;default:false" json:"visible_in_tab"`
	UpperLevelCategoryID *int32      `gorm:"type:int;default:null" json:"upper_level_category_id"`
	UpperLevelCategory   *Category   `gorm:"foreignKey:UpperLevelCategoryID;references:IncrementID" json:"upper_level_category"` // "belongs to" relationship with the `Category` model,
	SubCategories        []*Category `gorm:"foreignKey:UpperLevelCategoryID;references:IncrementID" json:"sub_categories"`       // "has many" relationship with the `Category` model,

}

type Brand struct {
	BaseModel
	Name string `gorm:"type:varchar(255);not null"`
	Logo string `gorm:"type:varchar(255);not null;default:''"`
}

type CategoryBrand struct {
	BaseModel
	CategoryID int32     `gorm:"type:int;index:idx_category_brand,unique"`
	Category   *Category // "belongs to" relationship with the `Category` model,
	BrandID    int32     `gorm:"type:int;index:idx_category_brand,unique"`
	Brand      *Brand    // "belongs to" relationship with the `Brand` model,

}

func (CategoryBrand) TableName() string {
	return "category_brand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(255);not null"`
	URL   string `gorm:"type:varchar(255);not null"`
	Index int32  `gorm:"type:int;not null;default:1"`
}

type Product struct {
	BaseModel
	CategoryID int32     `gorm:"type:int;not null"`
	Category   *Category // "belongs to" relationship with the `Category` model,
	BrandID    int32     `gorm:"type:int;not null"`
	Brand      *Brand    // "belongs to" relationship with the `Brand` model,

	OnSale       bool `gorm:"default:false;not null"`
	FreeShipping bool `gorm:"default:false;not null"`
	IsNew        bool `gorm:"default:false;not null"`
	IsHot        bool `gorm:"default:false;not null"`

	Name         string `gorm:"type:varchar(255);not null"`
	SerialNumber string `gorm:"type:varchar(50);not null"`
	Brief        string `gorm:"type:varchar(255);not null"`
	Description  string `gorm:"type:varchar(2000);not null"`

	ClickCount int32 `gorm:"type:int;not null;default:0"`
	LikeCount  int32 `gorm:"type:int;not null;default:0"`
	SoldCount  int32 `gorm:"type:int;not null;default:0"`

	MarketPrice float32 `gorm:"not null"`
	ShopPrice   float32 `gorm:"not null"`

	Images     StringList `gorm:"type:json;not null"`
	DescImages StringList `gorm:"type:json;not null"`
	FrontImage string     `gorm:"type:varchar(255);not null"`
}
