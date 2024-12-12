package main

import (
	"user_jwt/internal/interface/routes"
	"user_jwt/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// データベース接続の初期化
	config.ConnectDB()

	// ルートの設定
	routes.SetupRoutes(r)

	r.Run() // デフォルトで localhost:8080
}
