package controller

import (
	"codezard-pos/db"
	"codezard-pos/dto"
	"codezard-pos/model"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct{}

func (p Product) FindAll(ctx *gin.Context) {
	categoryId := ctx.Query("categoryId")
	search := ctx.Query("search")
	status := ctx.Query("status")
	// db.Conn.where("category_id=?").Find(&products)
	var products []model.Product
	// db.Conn.Find(&products) ดึงค่าตรงแต่จะไม่ได้ค่าความสัมพันธ์ที่อยู่ใน model_product แก้โดยใช้คำสั่งด้านล่าง
	// db.Conn.Preload("Category").Find(&products)
	query := db.Conn.Preload("Category")
	if categoryId != "" {
		query = query.Where("category_id=?", categoryId)
	}
	if search != "" {
		// ค้นหาแบบ like %-%
		query = query.Where("sku=? OR name LIKE ?", search, "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status=?", status)
	}

	query.Order("id desc").Find(&products)
	var result []dto.ReadProductResponse
	for _, product := range products {
		result = append(result, dto.ReadProductResponse{
			ID:     product.ID,
			SKU:    product.SKU,
			Name:   product.Name,
			Desc:   product.Desc,
			Price:  product.Price,
			Status: product.Status,
			Image:  product.Image,
			Category: dto.CategoryResponse{
				ID:   product.CategoryID,
				Name: product.Name,
			},
		})
	}
	ctx.JSON(http.StatusOK, result)
}

func (p Product) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")
	var product model.Product
	query := db.Conn.Preload("Category").First(&product, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ReadProductResponse{
		ID:     product.ID,
		SKU:    product.SKU,
		Name:   product.Name,
		Desc:   product.Desc,
		Price:  product.Price,
		Status: product.Status,
		Image:  product.Image,
		Category: dto.CategoryResponse{
			ID:   product.CategoryID,
			Name: product.Name,
		},
	})
}

func (p Product) Create(ctx *gin.Context) {

	var form dto.ProductRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	imagePath := "./uploads/products/" + uuid.New().String()
	ctx.SaveUploadedFile(image, imagePath)

	product := model.Product{
		SKU:        form.SKU,
		Name:       form.Name,
		Desc:       form.Desc,
		Price:      form.Price,
		Status:     form.Status,
		CategoryID: form.CategoryID,
		Image:      imagePath,
	}
	if err := db.Conn.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dto.CreateOrUpdateProductResponse{
		ID:         product.ID,
		SKU:        product.SKU,
		Name:       product.Name,
		Desc:       product.Desc,
		Price:      product.Price,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	})
}

func (p Product) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var form dto.ProductRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product model.Product
	query := db.Conn.First(&product, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	image, err := ctx.FormFile("image")
	// เพิ่มเงื่อนไขไม่ส่ง file ให้ผ่าน --> && !errors.Is(err, http.ErrMissingFile) ด้านล่างต่างจากการเช็ค error ปกติ

	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if image != nil {
		imagePath := "./uploads/products/" + uuid.New().String()
		ctx.SaveUploadedFile(image, imagePath)
		os.Remove(product.Image)
		product.Image = imagePath
	}
	product.SKU = form.SKU
	product.Name = form.Name
	product.Desc = form.Desc
	product.Price = form.Price
	product.Status = form.Status
	product.CategoryID = form.CategoryID
	db.Conn.Save(&product)
	ctx.JSON(http.StatusOK, dto.CreateOrUpdateProductResponse{
		ID:         product.ID,
		SKU:        product.SKU,
		Name:       product.Name,
		Desc:       product.Desc,
		Price:      product.Price,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	})
}

func (p Product) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	// unscoped() เพิ่มไปเพื่อลบออกจาก database จริงถ้าไม่ใส่่จะเป็นการ sofe_delete
	db.Conn.Unscoped().Delete(&model.Product{}, id)
	ctx.JSON(http.StatusOK, gin.H{"deletedAt": time.Now()})
}
