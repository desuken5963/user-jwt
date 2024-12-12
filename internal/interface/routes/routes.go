package routes

import (
	"user_jwt/internal/interface/handler"
	"user_jwt/internal/interface/repository"
	"user_jwt/internal/usecase"
	"user_jwt/pkg/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := config.DB

	// リポジトリ、ユースケース、ハンドラーを初期化
	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	// ルート設定
	auth := router.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/signin", authHandler.Signin)
	}
}
