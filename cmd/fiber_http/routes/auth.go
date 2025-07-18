package routes

import (
	"onfly-api/cmd/fiber_http/handlers"
	"onfly-api/internal/application/auth"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(router fiber.Router) {
	authService := auth.NewAuthService()
	authHandler := handlers.NewAuthHandler(authService)

	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)
	router.Post("/reset-password", authHandler.ResetPassword)
}
