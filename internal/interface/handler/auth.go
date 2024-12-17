package handler

import (
	"net/http"

	"user_jwt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// @Summary      Sign Up
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  struct{email string; password string; password_confirmation string}  true  "SignUp payload"
// @Success      201   {object} gin.H{"user": struct{ID uint; Email string}}
// @Failure      400   {object} gin.H{"error": string}
// @Failure      409   {object} gin.H{"error": string}
// @Router       /auth/sign-up [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req struct {
		Email                string `json:"email" binding:"required,email"`
		Password             string `json:"password" binding:"required,min=8"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}

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

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// @Summary      Sign In
// @Description  Authenticate a user and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  struct{email string; password string}  true  "SignIn payload"
// @Success      200   {object} gin.H{"token": string}
// @Failure      400   {object} gin.H{"error": string}
// @Failure      401   {object} gin.H{"error": string}
// @Router       /auth/sign-in [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	token, err := h.authUsecase.SignIn(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
