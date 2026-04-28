package handler

import (
	"Re_Shop/Backend/internal/modules/product/model"
	"Re_Shop/Backend/internal/modules/product/service"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) List(c *gin.Context) {
	query := buildQuery(c)
	products, err := h.productService.ListPublicProducts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "load products failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    products,
	})
}

func (h *ProductHandler) Detail(c *gin.Context) {
	productID, ok := parseID(c)
	if !ok {
		return
	}
	detail, err := h.productService.GetPublicProductDetail(productID)
	if err != nil {
		writeProductError(c, err, "load product detail failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    detail,
	})
}

func (h *ProductHandler) Categories(c *gin.Context) {
	categories, err := h.productService.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "load categories failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    categories,
	})
}

func (h *ProductHandler) Brands(c *gin.Context) {
	brands, err := h.productService.ListBrands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "load brands failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    brands,
	})
}

func (h *ProductHandler) AdminBrands(c *gin.Context) {
	brands, err := h.productService.ListAdminBrands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "load brands failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    brands,
	})
}

func (h *ProductHandler) AdminList(c *gin.Context) {
	query := buildQuery(c)
	products, err := h.productService.ListAdminProducts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "load admin products failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    products,
	})
}

func (h *ProductHandler) AdminDetail(c *gin.Context) {
	productID, ok := parseID(c)
	if !ok {
		return
	}
	detail, err := h.productService.GetAdminProductDetail(productID)
	if err != nil {
		writeProductError(c, err, "load admin product detail failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    detail,
	})
}

func (h *ProductHandler) AdminCreate(c *gin.Context) {
	var req model.ProductSaveInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "request payload is invalid",
		})
		return
	}

	detail, err := h.productService.CreateProduct(req)
	if err != nil {
		writeProductError(c, err, "create product failed")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "create product success",
		"data":    detail,
	})
}

func (h *ProductHandler) AdminUpdate(c *gin.Context) {
	productID, ok := parseID(c)
	if !ok {
		return
	}
	var req model.ProductSaveInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "request payload is invalid",
		})
		return
	}

	detail, err := h.productService.UpdateProduct(productID, req)
	if err != nil {
		writeProductError(c, err, "update product failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "update product success",
		"data":    detail,
	})
}

func (h *ProductHandler) AdminUpdateStatus(c *gin.Context) {
	productID, ok := parseID(c)
	if !ok {
		return
	}
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "request payload is invalid",
		})
		return
	}

	if err := h.productService.UpdateProductStatus(productID, req.Status); err != nil {
		writeProductError(c, err, "update product status failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "update product status success",
	})
}

func buildQuery(c *gin.Context) model.ProductQuery {
	query := model.ProductQuery{
		Keyword: strings.TrimSpace(c.Query("keyword")),
	}
	if categoryID, err := strconv.ParseInt(c.Query("categoryId"), 10, 64); err == nil {
		query.CategoryID = categoryID
	}
	if brandID, err := strconv.ParseInt(c.Query("brandId"), 10, 64); err == nil {
		query.BrandID = brandID
	}
	if statusText := c.Query("status"); statusText != "" {
		if parsedStatus, err := strconv.ParseInt(statusText, 10, 8); err == nil {
			status := int8(parsedStatus)
			query.Status = &status
		}
	}
	return query
}

func parseID(c *gin.Context) (int64, bool) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "product id is invalid",
		})
		return 0, false
	}
	return productID, true
}

func writeProductError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrInvalidProductPayload):
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	case errors.Is(err, service.ErrProductNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "product not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": fallback,
		})
	}
}
