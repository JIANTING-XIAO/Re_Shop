package service

import (
	"Re_Shop/Backend/internal/modules/product/model"
	"Re_Shop/Backend/internal/modules/product/repository"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

var (
	ErrInvalidProductPayload = errors.New("invalid product payload")
	ErrProductNotFound       = errors.New("product not found")
)

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) ListProducts() ([]model.Product, error) {
	return s.productRepo.ListEnabled()
}

func (s *ProductService) ListCategories() ([]model.Category, error) {
	return s.productRepo.ListCategories()
}

func (s *ProductService) ListBrands() ([]model.Brand, error) {
	return s.productRepo.ListBrands(true)
}

func (s *ProductService) ListAdminBrands() ([]model.Brand, error) {
	return s.productRepo.ListBrands(false)
}

func (s *ProductService) ListPublicProducts(query model.ProductQuery) ([]model.Product, error) {
	return s.productRepo.ListProducts(query, true)
}

func (s *ProductService) ListAdminProducts(query model.ProductQuery) ([]model.Product, error) {
	return s.productRepo.ListProducts(query, false)
}

func (s *ProductService) GetPublicProductDetail(productID int64) (*model.ProductDetail, error) {
	detail, err := s.productRepo.GetProductDetail(productID, true)
	return mapProductError(detail, err)
}

func (s *ProductService) GetAdminProductDetail(productID int64) (*model.ProductDetail, error) {
	detail, err := s.productRepo.GetProductDetail(productID, false)
	return mapProductError(detail, err)
}

func (s *ProductService) CreateProduct(input model.ProductSaveInput) (*model.ProductDetail, error) {
	if err := validateProductInput(input); err != nil {
		return nil, err
	}
	return s.productRepo.CreateProduct(input)
}

func (s *ProductService) UpdateProduct(productID int64, input model.ProductSaveInput) (*model.ProductDetail, error) {
	if err := validateProductInput(input); err != nil {
		return nil, err
	}
	detail, err := s.productRepo.UpdateProduct(productID, input)
	return mapProductError(detail, err)
}

func (s *ProductService) UpdateProductStatus(productID int64, status int8) error {
	if status != model.ProductStatusOnShelf && status != model.ProductStatusOffShelf {
		return fmt.Errorf("%w: status must be 0 or 1", ErrInvalidProductPayload)
	}
	err := s.productRepo.UpdateProductStatus(productID, status)
	_, mappedErr := mapProductError[struct{}](nil, err)
	return mappedErr
}

func validateProductInput(input model.ProductSaveInput) error {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.SpuCode) == "" {
		return fmt.Errorf("%w: name and spuCode are required", ErrInvalidProductPayload)
	}
	if input.CategoryID <= 0 || input.BrandID <= 0 {
		return fmt.Errorf("%w: categoryId and brandId are required", ErrInvalidProductPayload)
	}
	if input.Status != model.ProductStatusOnShelf && input.Status != model.ProductStatusOffShelf {
		return fmt.Errorf("%w: status must be 0 or 1", ErrInvalidProductPayload)
	}
	if len(input.SKUs) == 0 {
		return fmt.Errorf("%w: skus are required", ErrInvalidProductPayload)
	}
	for index, sku := range input.SKUs {
		if strings.TrimSpace(sku.SkuCode) == "" || strings.TrimSpace(sku.Name) == "" {
			return fmt.Errorf("%w: sku %d must have skuCode and name", ErrInvalidProductPayload, index+1)
		}
		if sku.Price <= 0 || sku.OriginalPrice < sku.Price {
			return fmt.Errorf("%w: sku %d price is invalid", ErrInvalidProductPayload, index+1)
		}
		if sku.Stock < 0 {
			return fmt.Errorf("%w: sku %d stock is invalid", ErrInvalidProductPayload, index+1)
		}
		if sku.Status != model.ProductStatusOnShelf && sku.Status != model.ProductStatusOffShelf {
			return fmt.Errorf("%w: sku %d status must be 0 or 1", ErrInvalidProductPayload, index+1)
		}
	}
	return nil
}

func mapProductError[T any](value *T, err error) (*T, error) {
	if err == nil {
		return value, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}
	return nil, err
}
