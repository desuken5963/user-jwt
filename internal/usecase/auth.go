package usecase

import (
	"errors"
	"user_jwt/internal/domain"
	"user_jwt/internal/repository"
	"user_jwt/pkg/utils"
)

// AuthUsecase インターフェース
type AuthUsecase interface {
	SignUp(email, password string) (domain.User, error)
	SignIn(email, password string) (string, error) // JWTトークンを返す
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) SignUp(email, password string) (domain.User, error) {
	// 重複チェック
	existingUser, _ := u.userRepo.FindByEmail(email)
	if existingUser != nil {
		return domain.User{}, errors.New("email already exists")
	}

	// パスワードハッシュ化
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	// ユーザー作成
	user := domain.User{
		Email:    email,
		Password: hashedPassword,
	}
	createdUser, err := u.userRepo.Create(user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (u *authUsecase) SignIn(email, password string) (string, error) {
	// ユーザー取得
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	// パスワードチェック
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	// JWTトークン生成
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
