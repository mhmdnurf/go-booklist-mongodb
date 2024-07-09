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
	"go.mongodb.org/mongo-driver/bson"
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

	if book.Author == "" || book.Title == " " || book.Year == 0 {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Data can't be empty"}})

	}

	newBook := models.Book{
		Id:        primitive.NewObjectID(),
		Title:     book.Title,
		Author:    book.Author,
		Year:      book.Year,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := bookCollection.InsertOne(ctx, newBook)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.BookResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"book": newBook, "inserted_id": result.InsertedID}})
}

func UpdateBook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var bookUpdate models.BookUpdate
	if err := c.BodyParser(&bookUpdate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if validationErr := validate.Struct(&bookUpdate); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{
			Status:  http.StatusBadRequest,
			Message: validationErr.Error(),
			Data:    nil,
		})
	}

	bookId := c.Params("id")
	id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	updateData := bson.D{}

	if bookUpdate.Title != nil {
		updateData = append(updateData, bson.E{Key: "title", Value: *bookUpdate.Title})
	}
	if bookUpdate.Author != nil {
		updateData = append(updateData, bson.E{Key: "author", Value: *bookUpdate.Author})
	}
	if bookUpdate.Year != nil {
		updateData = append(updateData, bson.E{Key: "year", Value: *bookUpdate.Year})
	}

	updateData = append(updateData, bson.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())})

	update := bson.D{{Key: "$set", Value: updateData}}

	_, updateErr := bookCollection.UpdateOne(ctx, primitive.D{{Key: "_id", Value: id}}, update)

	if updateErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{
			Status:  http.StatusInternalServerError,
			Message: updateErr.Error(),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responses.BookResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    &fiber.Map{"book": bookUpdate},
	})
}

func DeleteBook(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookId := c.Params("id")
	id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BookResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	_, deleteErr := bookCollection.DeleteOne(ctx, primitive.D{{Key: "_id", Value: id}})
	if deleteErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BookResponse{
			Status:  http.StatusInternalServerError,
			Message: deleteErr.Error(),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responses.BookResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}
