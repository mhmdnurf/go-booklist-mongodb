package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mhmdnurf/go-booklist-mongodb/configs"
	"github.com/mhmdnurf/go-booklist-mongodb/models"
	"github.com/mhmdnurf/go-booklist-mongodb/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookCollection *mongo.Collection = configs.GetCollection(configs.DB, "books")
var validate = validator.New()

func GetAllBooks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bookCollection.Find(ctx, primitive.D{{}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var book []models.Book = make([]models.Book, 0)

	for cursor.Next(ctx) {
		var b models.Book
		cursor.Decode(&b)
		book = append(book, b)
	}
	return c.Status(http.StatusOK).JSON(responses.BookResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"books": book}})
}

func CreateBook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var book models.Book
	defer cancel()

	if err := c.BodyParser(&book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&book); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newBook := models.Book{
		Id:        primitive.NewObjectID(),
		Title:     book.Title,
		Author:    book.Author,
		Year:      book.Year,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := bookCollection.InsertOne(ctx, newBook)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.BookResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"book": newBook, "inserted_id": result.InsertedID}})
}
