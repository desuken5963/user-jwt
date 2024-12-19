package repository

import (
	"user_jwt/internal/domain"
)

// UserRepository インターフェース
type UserRepository interface {
	FindByEmail(email string) (*domain.User, error) // ユーザーをメールで検索
	Create(user domain.User) (domain.User, error)   // ユーザーを作成
	FindByID(userID uint) (*domain.User, error)
}
