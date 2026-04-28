package repository

import (
	"Re_Shop/Backend/internal/modules/product/model"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) ListEnabled() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("status = ?", 1).Order("id asc").Find(&products).Error
	return products, err
}

func (r *ProductRepository) ListCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort asc, id asc").Find(&categories).Error
	return categories, err
}

func (r *ProductRepository) ListBrands(enabledOnly bool) ([]model.Brand, error) {
	var brands []model.Brand
	query := r.db.Order("sort asc, id asc")
	if enabledOnly {
		query = query.Where("status = ?", model.ProductStatusOnShelf)
	}
	err := query.Find(&brands).Error
	return brands, err
}

func (r *ProductRepository) ListProducts(query model.ProductQuery, enabledOnly bool) ([]model.Product, error) {
	var products []model.Product
	dbQuery := r.db.Model(&model.Product{}).Order("id desc")
	if enabledOnly {
		dbQuery = dbQuery.Where("status = ?", model.ProductStatusOnShelf)
	} else if query.Status != nil {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}
	if query.Keyword != "" {
		like := "%" + strings.TrimSpace(query.Keyword) + "%"
		dbQuery = dbQuery.Where("name LIKE ? OR title LIKE ? OR spu_code LIKE ?", like, like, like)
	}
	if query.CategoryID > 0 {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
	}
	if query.BrandID > 0 {
		dbQuery = dbQuery.Where("brand_id = ?", query.BrandID)
	}
	err := dbQuery.Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetProductDetail(productID int64, enabledOnly bool) (*model.ProductDetail, error) {
	var product model.Product
	query := r.db.Where("id = ?", productID)
	if enabledOnly {
		query = query.Where("status = ?", model.ProductStatusOnShelf)
	}
	if err := query.First(&product).Error; err != nil {
		return nil, err
	}

	var category model.Category
	if err := r.db.Where("id = ?", product.CategoryID).First(&category).Error; err != nil {
		return nil, err
	}

	var brand model.Brand
	if err := r.db.Where("id = ?", product.BrandID).First(&brand).Error; err != nil {
		return nil, err
	}

	var skus []model.ProductSKU
	skuQuery := r.db.Where("product_id = ?", product.ID).Order("id asc")
	if enabledOnly {
		skuQuery = skuQuery.Where("status = ?", model.ProductStatusOnShelf)
	}
	if err := skuQuery.Find(&skus).Error; err != nil {
		return nil, err
	}

	return &model.ProductDetail{
		Product:  product,
		Category: category,
		Brand:    brand,
		SKUs:     skus,
	}, nil
}

func (r *ProductRepository) CreateProduct(input model.ProductSaveInput) (*model.ProductDetail, error) {
	var detail *model.ProductDetail
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureCategoryExists(tx, input.CategoryID); err != nil {
			return err
		}
		if err := ensureBrandExists(tx, input.BrandID); err != nil {
			return err
		}

		price, originalPrice, totalStock, err := summarizeSKUs(input.SKUs)
		if err != nil {
			return err
		}

		product := model.Product{
			Name:          input.Name,
			Title:         input.Title,
			SpuCode:       input.SpuCode,
			CategoryID:    input.CategoryID,
			BrandID:       input.BrandID,
			Price:         price,
			OriginalPrice: originalPrice,
			CoverImage:    input.CoverImage,
			Detail:        input.Detail,
			Status:        input.Status,
		}
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		skus, err := buildSKUs(product.ID, input.SKUs)
		if err != nil {
			return err
		}
		if err := tx.Create(&skus).Error; err != nil {
			return err
		}

		stock := model.ProductStock{
			ProductID:   product.ID,
			Stock:       totalStock,
			LockedStock: 0,
		}
		if err := tx.Where("product_id = ?", product.ID).Assign(stock).FirstOrCreate(&stock).Error; err != nil {
			return err
		}

		detail, err = loadProductDetail(tx, product.ID, false)
		return err
	})
	return detail, err
}

func (r *ProductRepository) UpdateProduct(productID int64, input model.ProductSaveInput) (*model.ProductDetail, error) {
	var detail *model.ProductDetail
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing model.Product
		if err := tx.Where("id = ?", productID).First(&existing).Error; err != nil {
			return err
		}
		if err := ensureCategoryExists(tx, input.CategoryID); err != nil {
			return err
		}
		if err := ensureBrandExists(tx, input.BrandID); err != nil {
			return err
		}

		price, originalPrice, totalStock, err := summarizeSKUs(input.SKUs)
		if err != nil {
			return err
		}

		updates := map[string]any{
			"name":           input.Name,
			"title":          input.Title,
			"spu_code":       input.SpuCode,
			"category_id":    input.CategoryID,
			"brand_id":       input.BrandID,
			"price":          price,
			"original_price": originalPrice,
			"cover_image":    input.CoverImage,
			"detail":         input.Detail,
			"status":         input.Status,
		}
		if err := tx.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}

		if err := tx.Where("product_id = ?", productID).Delete(&model.ProductSKU{}).Error; err != nil {
			return err
		}
		skus, err := buildSKUs(productID, input.SKUs)
		if err != nil {
			return err
		}
		if err := tx.Create(&skus).Error; err != nil {
			return err
		}

		stock := model.ProductStock{
			ProductID:   productID,
			Stock:       totalStock,
			LockedStock: 0,
		}
		if err := tx.Where("product_id = ?", productID).Assign(stock).FirstOrCreate(&stock).Error; err != nil {
			return err
		}

		detail, err = loadProductDetail(tx, productID, false)
		return err
	})
	return detail, err
}

func (r *ProductRepository) UpdateProductStatus(productID int64, status int8) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Product{}).
			Where("id = ?", productID).
			Updates(map[string]any{"status": status})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&model.ProductSKU{}).
			Where("product_id = ?", productID).
			Update("status", status).Error
	})
}

func ensureCategoryExists(tx *gorm.DB, categoryID int64) error {
	var count int64
	if err := tx.Model(&model.Category{}).Where("id = ?", categoryID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("category %d not found", categoryID)
	}
	return nil
}

func ensureBrandExists(tx *gorm.DB, brandID int64) error {
	var count int64
	if err := tx.Model(&model.Brand{}).Where("id = ?", brandID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("brand %d not found", brandID)
	}
	return nil
}

func summarizeSKUs(inputs []model.ProductSKUInput) (float64, float64, int, error) {
	if len(inputs) == 0 {
		return 0, 0, 0, errors.New("at least one sku is required")
	}
	minPrice := inputs[0].Price
	minOriginalPrice := inputs[0].OriginalPrice
	totalStock := 0
	for _, sku := range inputs {
		if sku.Price < minPrice {
			minPrice = sku.Price
		}
		if sku.OriginalPrice < minOriginalPrice {
			minOriginalPrice = sku.OriginalPrice
		}
		totalStock += sku.Stock
	}
	return minPrice, minOriginalPrice, totalStock, nil
}

func buildSKUs(productID int64, inputs []model.ProductSKUInput) ([]model.ProductSKU, error) {
	skus := make([]model.ProductSKU, 0, len(inputs))
	for _, input := range inputs {
		specs := "{}"
		if len(input.Specs) > 0 {
			data, err := json.Marshal(input.Specs)
			if err != nil {
				return nil, err
			}
			specs = string(data)
		}
		skus = append(skus, model.ProductSKU{
			ProductID:     productID,
			SkuCode:       input.SkuCode,
			Name:          input.Name,
			Specs:         specs,
			Price:         input.Price,
			OriginalPrice: input.OriginalPrice,
			Stock:         input.Stock,
			Status:        input.Status,
		})
	}
	return skus, nil
}

func loadProductDetail(tx *gorm.DB, productID int64, enabledOnly bool) (*model.ProductDetail, error) {
	repo := &ProductRepository{db: tx}
	return repo.GetProductDetail(productID, enabledOnly)
}
