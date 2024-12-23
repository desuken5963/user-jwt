package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"user-jwt/pkg/utils"
)

// AuthMiddleware JWTトークンを検証するミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーからトークンを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// トークンの前に "Bearer " が付いている場合を考慮
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// トークンを検証
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 検証成功後、コンテキストにユーザー情報を保存
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		// 次の処理に進む
		c.Next()
	}
}
