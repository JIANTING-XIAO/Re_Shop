package db

import (
	"fmt"

	productModel "Re_Shop/Backend/internal/modules/product/model"
	userModel "Re_Shop/Backend/internal/modules/user/model"
	"Re_Shop/Backend/internal/shared/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	dbCfg := cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DBName,
		dbCfg.Charset,
		dbCfg.ParseTime,
		dbCfg.Loc,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	if err := DB.AutoMigrate(
		&userModel.User{},
		&productModel.Category{},
		&productModel.Brand{},
		&productModel.Product{},
		&productModel.ProductSKU{},
		&productModel.ProductStock{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate tables: %w", err)
	}

	return DB, nil
}

func Get() *gorm.DB {
	return DB
}
