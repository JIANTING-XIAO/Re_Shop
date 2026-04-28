package main

import (
	"Re_Shop/Backend/internal/app/router"
	"Re_Shop/Backend/internal/shared/db"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	if _, err := db.Init(); err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	r := gin.Default()
	registerStaticAssets(r)
	registerFrontendPage(r)
	router.Register(r)
	router.RegisterUserRoutes(r)

	log.Println("server running at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// 加载静态资源
func registerStaticAssets(r *gin.Engine) {
	assetPath, err := resolveAssetPath()
	if err != nil {
		log.Fatalf("resolve static asset path failed: %v", err)
	}

	r.Static("/static/img", assetPath)
}

func resolveAssetPath() (string, error) {
	candidates := []string{
		filepath.Join("Backend", "internal", "resource", "img"),
		filepath.Join("internal", "resource", "img"),
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", os.ErrNotExist
}

func registerFrontendPage(r *gin.Engine) {
	indexPath, err := resolveFrontendIndexPath()
	if err != nil {
		log.Printf("frontend index not found: %v", err)
		return
	}

	r.GET("/index", func(c *gin.Context) {
		c.File(indexPath)
	})
}

func resolveFrontendIndexPath() (string, error) {
	candidates := []string{
		filepath.Join("Frontend", "index.html"),
		filepath.Join("..", "Frontend", "index.html"),
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", os.ErrNotExist
}
