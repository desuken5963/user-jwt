package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTの署名キー
var jwtKey = []byte("your_secret_key")

// カスタムクレーム
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTトークンを生成
func GenerateJWT(userID uint, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24時間の有効期限

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// JWTトークンを検証
func VerifyJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
