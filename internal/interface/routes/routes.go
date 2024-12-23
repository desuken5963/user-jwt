package routes

import (
	"user-jwt/internal/interface/handler"
	"user-jwt/internal/interface/middleware"
	"user-jwt/internal/interface/repository"
	"user-jwt/internal/usecase"
	"user-jwt/pkg/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := config.DB

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/sign-in", authHandler.SignIn)
	}

	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/:id", userHandler.GetUserByID)
	}
}
