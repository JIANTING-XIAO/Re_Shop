package model

import "time"

const (
	ProductStatusOffShelf int8 = 0
	ProductStatusOnShelf  int8 = 1
)

type Category struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	ParentID  int64     `gorm:"column:parent_id;not null;default:0" json:"parentId"`
	Sort      int       `gorm:"column:sort;not null;default:0" json:"sort"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Category) TableName() string {
	return "categories"
}

type Brand struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Logo        string    `gorm:"column:logo;type:varchar(500)" json:"logo"`
	Description string    `gorm:"column:description;type:varchar(255)" json:"description"`
	Sort        int       `gorm:"column:sort;not null;default:0" json:"sort"`
	Status      int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Brand) TableName() string {
	return "brands"
}

type Product struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"column:name;type:varchar(200);not null" json:"name"`
	Title         string    `gorm:"column:title;type:varchar(255)" json:"title"`
	SpuCode       string    `gorm:"column:spu_code;type:varchar(64);not null" json:"spuCode"`
	CategoryID    int64     `gorm:"column:category_id;not null" json:"categoryId"`
	BrandID       int64     `gorm:"column:brand_id;not null;default:0" json:"brandId"`
	Price         float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64   `gorm:"column:original_price;type:decimal(10,2);not null" json:"originalPrice"`
	CoverImage    string    `gorm:"column:cover_image;type:varchar(500)" json:"coverImage"`
	Detail        string    `gorm:"column:detail;type:text" json:"detail"`
	Status        int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Product) TableName() string {
	return "products"
}

type ProductSKU struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID     int64     `gorm:"column:product_id;not null" json:"productId"`
	SkuCode       string    `gorm:"column:sku_code;type:varchar(64);not null" json:"skuCode"`
	Name          string    `gorm:"column:name;type:varchar(200);not null" json:"name"`
	Specs         string    `gorm:"column:specs;type:json" json:"specs"`
	Price         float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64   `gorm:"column:original_price;type:decimal(10,2);not null" json:"originalPrice"`
	Stock         int       `gorm:"column:stock;not null;default:0" json:"stock"`
	Status        int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ProductSKU) TableName() string {
	return "product_skus"
}

type ProductStock struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ProductID   int64     `gorm:"column:product_id;not null;uniqueIndex"`
	Stock       int       `gorm:"column:stock;not null;default:0"`
	LockedStock int       `gorm:"column:locked_stock;not null;default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ProductStock) TableName() string {
	return "product_stock"
}

type ProductQuery struct {
	Keyword    string
	CategoryID int64
	BrandID    int64
	Status     *int8
}

type ProductSKUInput struct {
	ID            int64             `json:"id"`
	SkuCode       string            `json:"skuCode"`
	Name          string            `json:"name"`
	Specs         map[string]string `json:"specs"`
	Price         float64           `json:"price"`
	OriginalPrice float64           `json:"originalPrice"`
	Stock         int               `json:"stock"`
	Status        int8              `json:"status"`
}

type ProductSaveInput struct {
	Name       string            `json:"name"`
	Title      string            `json:"title"`
	SpuCode    string            `json:"spuCode"`
	CategoryID int64             `json:"categoryId"`
	BrandID    int64             `json:"brandId"`
	CoverImage string            `json:"coverImage"`
	Detail     string            `json:"detail"`
	Status     int8              `json:"status"`
	SKUs       []ProductSKUInput `json:"skus"`
}

type ProductDetail struct {
	Product  Product      `json:"product"`
	Category Category     `json:"category"`
	Brand    Brand        `json:"brand"`
	SKUs     []ProductSKU `json:"skus"`
}
