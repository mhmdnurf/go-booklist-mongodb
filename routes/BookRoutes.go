package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhmdnurf/go-booklist-mongodb/controllers"
)

func BookRoute(app *fiber.App) {
	app.Get("/books", controllers.GetAllBooks)
	app.Post("/book", controllers.CreateBook)
	app.Put("/book/:id", controllers.UpdateBook)
	app.Delete("/book/:id", controllers.DeleteBook)
}
