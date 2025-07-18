package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"onfly-api/cmd/fiber_http/routes"
	_ "onfly-api/docs"
	"onfly-api/internal/infrasctructure/persistence"

	swagger "github.com/gofiber/swagger"
)

// @title Onfly API
// @version 1.0
// @description Esta é a documentação da API da Onfly
// @host localhost:3000
// @BasePath /
// @schemes http

func StartServerHttp() {
	persistence.InitDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.RegisterRoutes(app)

	app.Listen(":3000")
}
