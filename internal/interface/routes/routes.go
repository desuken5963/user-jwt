package routes

import (
	"user_jwt/internal/interface/handler"
	"user_jwt/internal/interface/middleware"
	"user_jwt/internal/interface/repository"
	"user_jwt/internal/usecase"
	"user_jwt/pkg/config"

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
