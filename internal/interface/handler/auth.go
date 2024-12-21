package handler

import (
	"net/http"

	"user-jwt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// リクエスト・レスポンス用構造体定義
type SignUpRequest struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type SignUpResponse struct {
	User UserResponse `json:"user"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

	token, err := h.authUsecase.SignIn(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := SignInResponse{Token: token}
	c.JSON(http.StatusOK, response)
}
