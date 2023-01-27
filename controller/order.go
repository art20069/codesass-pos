package controller

import (
	"codezard-pos/db"
	"codezard-pos/dto"
	"codezard-pos/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Order struct{}

func (o Order) FindAll(ctx *gin.Context) {
	var orders []model.Order
	db.Conn.Preload("Products").Find(&orders)

	var result []dto.OrderResponse
	for _, order := range orders {
		orderResult := dto.OrderResponse{
			ID:    order.ID,
			Name:  order.Name,
			Tel:   order.Tel,
			Email: order.Email,
		}
		var products []dto.OrderProductResponse
		for _, product := range order.Products {
			products = append(products, dto.OrderProductResponse{
				ID:       product.ID,
				SKU:      product.SKU,
				Name:     product.Name,
				Price:    product.Price,
				Quantity: product.Quantity,
				Image:    product.Image,
			})
		}
		orderResult.Products = products
		result = append(result, orderResult)
	}

	ctx.JSON(http.StatusOK, result)
}

func (o Order) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")
	var order model.Order

	query := db.Conn.Preload("Products").First(&order, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := dto.OrderResponse{
		ID:    order.ID,
		Name:  order.Name,
		Tel:   order.Tel,
		Email: order.Email,
	}
	var products []dto.OrderProductResponse
	for _, product := range order.Products {
		products = append(products, dto.OrderProductResponse{
			ID:       product.ID,
			SKU:      product.SKU,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quantity,
			Image:    product.Image,
		})
	}
	result.Products = products

	ctx.JSON(http.StatusOK, result)
}

func (o Order) Create(ctx *gin.Context) {
	var form dto.OrderRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order model.Order
	var orderItems []model.OrderItem
	for _, product := range form.Products {
		orderItems = append(orderItems, model.OrderItem{
			SKU:      product.SKU,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quantity,
			Image:    product.Image,
		})
	}
	order.Name = form.Name
	order.Tel = form.Tel
	order.Email = form.Email
	order.Products = orderItems
	db.Conn.Create(&order)

	result := dto.OrderResponse{
		ID:    order.ID,
		Name:  order.Name,
		Tel:   order.Tel,
		Email: order.Email,
	}

	var products []dto.OrderProductResponse

	for _, product := range order.Products {
		products = append(products, dto.OrderProductResponse{
			ID:       product.ID,
			SKU:      product.SKU,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quantity,
			Image:    product.Image,
		})
	}

	result.Products = products
	// result.Products = products

	ctx.JSON(http.StatusCreated, result)
}
