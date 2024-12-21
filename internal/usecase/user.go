package usecase

import (
	"errors"
	"user-jwt/internal/domain"
	"user-jwt/internal/repository"
)

// UserUsecase ユーザーに関するユースケース
type UserUsecase interface {
	GetUserByID(userID uint) (*domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase UserUsecaseのコンストラクタ
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

// GetUserByID ユーザーIDでユーザー情報を取得
func (u *userUsecase) GetUserByID(userID uint) (*domain.User, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
