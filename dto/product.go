package dto

type ProductRequest struct {
	Name       string  `form:"name" binding:"required"`
	SKU        string  `form:"sku" binding:"required"`
	Desc       string  `form:"desc"`
	Price      float64 `form:"price" binding:"required"`
	Status     uint    `form:"status" binding:"required"`
	CategoryID uint    `form:"categoryId" binding:"required"`
}

type CreateOrUpdateProductResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Desc       string  `json:"desc"`
	Price      float64 `json:"price"`
	Status     uint    `json:"status"`
	CategoryID uint    `json:"categoryId"`
	Image      string  `json:"image"`
}

type ReadProductResponse struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	SKU      string           `json:"sku"`
	Desc     string           `json:"desc"`
	Price    float64          `json:"price"`
	Status   uint             `json:"status"`
	Category CategoryResponse `json:"category"` // เปลี่ยนแปลงค่าโดยนำค่าจาก dto.category.go มาใช้งานแทนรูปแบบ ID เดิม
	Image    string           `json:"image"`
}
