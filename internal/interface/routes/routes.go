package routes

import (
	"net/http"

	"user_jwt/internal/interface/handler"
	"user_jwt/internal/interface/middleware"
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
		auth.POST("/sign-up", authHandler.Signup)
		auth.POST("/sign-in", authHandler.Signin)
	}

	// 認証が必要なエンドポイント
	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware()) // ミドルウェアを適用
	{
		protected.GET("/user-info", func(c *gin.Context) {
			// コンテキストからユーザー情報を取得
			userID := c.GetInt("userID")
			email := c.GetString("email")

			c.JSON(http.StatusOK, gin.H{"userID": userID, "email": email})
		})
	}
}
