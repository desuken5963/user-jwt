package main

import (
	_ "user_jwt/docs"
	"user_jwt/internal/interface/routes"
	"user_jwt/pkg/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Swaggerエンドポイント
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// データベース接続の初期化
	config.ConnectDB()

	// ルートの設定
	routes.SetupRoutes(r)

	r.Run() // デフォルトで localhost:8080
}
