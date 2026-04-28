package router

import (
	"Re_Shop/Backend/internal/modules/user/handler"
	"Re_Shop/Backend/internal/modules/user/repository"
	"Re_Shop/Backend/internal/modules/user/service"
	"Re_Shop/Backend/internal/shared/db"
	"Re_Shop/Backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userRepo := repository.NewUserRepository(db.Get())
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)
	r.GET("/me", middleware.AuthRequired(), authHandler.Me)
}
