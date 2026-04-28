package router

import (
	"Re_Shop/Backend/internal/modules/product/handler"
	"Re_Shop/Backend/internal/modules/product/repository"
	"Re_Shop/Backend/internal/modules/product/service"
	"Re_Shop/Backend/internal/shared/db"
	"Re_Shop/Backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine) {
	productRepo := repository.NewProductRepository(db.Get())
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r.GET("/products", productHandler.List)
	r.GET("/products/:id", productHandler.Detail)
	r.GET("/categories", productHandler.Categories)
	r.GET("/brands", productHandler.Brands)

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired())
	adminGroup.GET("/brands", productHandler.AdminBrands)
	adminGroup.GET("/products", productHandler.AdminList)
	adminGroup.GET("/products/:id", productHandler.AdminDetail)
	adminGroup.POST("/products", productHandler.AdminCreate)
	adminGroup.PUT("/products/:id", productHandler.AdminUpdate)
	adminGroup.PATCH("/products/:id/status", productHandler.AdminUpdateStatus)
}
