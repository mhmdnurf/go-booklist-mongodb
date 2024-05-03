package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhmdnurf/go-booklist-mongodb/configs"
	"github.com/mhmdnurf/go-booklist-mongodb/routes"
)

func main() {
	app := fiber.New()
	configs.ConnectDB()
	routes.BookRoute(app)
	app.Listen(":3000")
}
