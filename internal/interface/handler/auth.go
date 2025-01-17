package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"user-jwt/internal/usecase"
	"user-jwt/pkg/config"
	"user-jwt/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// SugnUPリクエスト・レスポンス用構造体定義
type SignUpRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type SignUpResponse struct {
	User UserResponse `json:"user"`
}

// SugnINリクエスト・レスポンス用構造体定義
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

// @Summary      Sign Up
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  SignUpRequest  true  "SignUp payload"
// @Success      201   {object} SignUpResponse
// @Failure      400   {object} map[string]string
// @Failure      409   {object} map[string]string
// @Router       /auth/sign-up [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	validationErrors := utils.ValidateStruct(&req)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	if req.Password != req.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password and confirmation do not match"})
		return
	}

	user, err := h.authUsecase.SignUp(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	response := SignUpResponse{
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary      Sign In
// @Description  Authenticate a user and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  SignInRequest  true  "SignIn payload"
// @Success      200   {object} SignInResponse
// @Failure      400   {object} map[string]string
// @Failure      401   {object} map[string]string
// @Router       /auth/sign-in [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	token, err := h.authUsecase.SignIn(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := SignInResponse{Token: token}
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) PostAuthSignOut(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	claims, err := utils.VerifyJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// トークンをRedisに追加
	expiration := time.Until(claims.ExpiresAt.Time)
	err = config.RedisClient.Set(context.Background(), tokenString, "revoked", expiration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed out"})
}
