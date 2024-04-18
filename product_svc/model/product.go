package model

type Category struct {
	BaseModel
	Name                 string `gorm:"type:varchar(20);not null"`
	Level                int32  `gorm:"type:int;not null;default:1"`
	VisibleInTab         bool   `gorm:"not null;default:false"`
	UpperLevelCategoryID *int32
	UpperLevelCategory   *Category
	SubCategories        []*Category `gorm:"foreignKey:UpperLevelCategoryID;references:ID"`
}

type Brand struct {
	BaseModel
	Name string `gorm:"type:varchar(255);not null"`
	Logo string `gorm:"type:varchar(255);not null;default:''"`
}

type CategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   *Category
	BrandID    int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brand      *Brand
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
	CategoryID int32 `gorm:"type:int;not null"`
	Category   *Category
	BrandID    int32 `gorm:"type:int;not null"`
	Brand      *Brand

	OnSale       bool `gorm:"default:false;not null"`
	FreeShipping bool `gorm:"default:false;not null"`
	IsNew        bool `gorm:"default:false;not null"`
	IsHot        bool `gorm:"default:false;not null"`

	Name         string `gorm:"type:varchar(255);not null"`
	SerialNumber string `gorm:"type:varchar(50);not null"`
	Brief        string `gorm:"type:varchar(255);not null"`

	ClickCount int32 `gorm:"type:int;not null;default:0"`
	LikeCount  int32 `gorm:"type:int;not null;default:0"`
	SoldCount  int32 `gorm:"type:int;not null;default:0"`

	MarketPrice float32 `gorm:"not null"`
	ShopPrice   float32 `gorm:"not null"`

	Images      StringList `gorm:"type:json;not null"`
	DescImages  StringList `gorm:"type:json;not null"`
	FrontImages string     `gorm:"type:varchar(255);not null"`
}
