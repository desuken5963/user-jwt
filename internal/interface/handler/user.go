package handler

import (
	"net/http"
	"strconv"

	"user-jwt/internal/usecase"
	"user-jwt/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

type GetUserByIDRequest struct {
	ID uint `validate:"required,min=1"`
}

// GetUserByID ユーザーIDで情報を取得
// @Summary      Get User by ID
// @Description  Retrieve user information using the user ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int     true  "User ID" Format(int64)
// @Success      200  {object}  map[string]interface{}  "Successful Response"
// @Failure      400  {object}  map[string]string       "Invalid User ID"
// @Failure      404  {object}  map[string]string       "User Not Found"
// @Router       /user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	req := GetUserByIDRequest{ID: uint(userID)}
	validationErrors := utils.ValidateStruct(&req)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	user, err := h.userUsecase.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userID": user.ID, "email": user.Email})
}
