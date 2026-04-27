package main

import (
	"Re_Shop/Backend/internal/app/router"
	"Re_Shop/Backend/internal/shared/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	//加载
	r := gin.Default()
	router.Register(r)
	router.RegesiterUserRoutes(r)

	log.Println("server running at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	//mysql
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}

	fmt.Println("database connected:", db != nil)
}
