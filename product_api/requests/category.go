package requests

type StoreCategoryRequest struct {
	Name                 string `json:"name" binding:"required,min=1,max=255"`
	UpperLevelCategoryId int32  `json:"upper_level_category_id" binding:"min=0"`
	IsTab                bool   `json:"is_tab"`
}

type UpdateCategoryRequest struct {
	Name                 string `json:"name" binding:"required,min=1,max=255"`
	UpperLevelCategoryId int32  `json:"upper_level_category_id" binding:"min=0"`
	IsTab                bool   `json:"is_tab"`
}
