package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhmdnurf/go-booklist-mongodb/controllers"
)

func BookRoute(app *fiber.App) {
	app.Get("/books", controllers.GetAllBooks)
	app.Post("/books", controllers.CreateBook)
}
