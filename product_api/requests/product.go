package requests

type StoreProductRequest struct {
	Name         string   `form:"name" json:"name" binding:"required,min=3,max=100"`
	CategoryID   int32    `form:"category_id" json:"category_id" binding:"required"`
	BrandID      int32    `form:"brand_id" json:"brand_id" binding:"required"`
	SerialNumber string   `form:"serial_number" json:"serial_number" binding:"required"`
	MarketPrice  float32  `form:"market_price" json:"market_price" binding:"required,min=0"`
	ShopPrice    float32  `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	Brief        string   `form:"brief" json:"brief" binding:"required"`
	Description  string   `form:"description" json:"description" binding:"required,min=3"`
	FreeShipping bool     `form:"free_shipping" json:"free_shipping" binding:"required"`
	Images       []string `form:"images" json:"images" binding:"required"`
	DescImages   []string `form:"desc_images" json:"desc_images" binding:"required"`
	FrontImage   string   `form:"front_image" json:"front_image" binding:"required"`
	Stock        int32    `form:"stock" json:"stock" binding:"required,min=1"`
}

type UpdateProductRequest struct {
	Name         string   `form:"name" json:"name" binding:"required,min=3,max=100"`
	CategoryID   int32    `form:"category_id" json:"category_id" binding:"required"`
	BrandID      int32    `form:"brand_id" json:"brand_id" binding:"required"`
	SerialNumber string   `form:"serial_number" json:"serial_number" binding:"required"`
	MarketPrice  float32  `form:"market_price" json:"market_price" binding:"required,min=0"`
	ShopPrice    float32  `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	Brief        string   `form:"brief" json:"brief" binding:"required"`
	Description  string   `form:"description" json:"description" binding:"required,min=3"`
	FreeShipping bool     `form:"free_shipping" json:"free_shipping" binding:"required"`
	Images       []string `form:"images" json:"images" binding:"required"`
	DescImages   []string `form:"desc_images" json:"desc_images" binding:"required"`
	FrontImage   string   `form:"front_image" json:"front_image" binding:"required"`
	Stock        int32    `form:"stock" json:"stock" binding:"required,min=1"`
	OnSale       bool     `form:"on_sale" json:"on_sale" binding:"required"`
	IsNew        bool     `form:"is_new" json:"is_new" binding:"required"`
	IsHot        bool     `form:"is_hot" json:"is_hot" binding:"required"`
}
