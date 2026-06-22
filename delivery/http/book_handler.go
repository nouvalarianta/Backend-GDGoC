package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nouval/library-app/delivery/http/middleware"
	"github.com/nouval/library-app/domain"
)

// BookHandler represent the httphandler for book
type BookHandler struct {
	BUsecase domain.BookUsecase
}

// NewBookHandler will initialize the books/ resources endpoint
func NewBookHandler(app *fiber.App, us domain.BookUsecase, jwtSecret string) {
	handler := &BookHandler{
		BUsecase: us,
	}

	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	app.Get("/books", handler.Fetch)
	app.Get("/books/:id", handler.GetByID)

	// Protected endpoints
	app.Post("/books", authMiddleware, handler.Store)
	app.Put("/books/:id", authMiddleware, handler.Update)
	app.Delete("/books/:id", authMiddleware, handler.Delete)
}

func (a *BookHandler) Fetch(c *fiber.Ctx) error {
	ctx := c.Context()
	books, err := a.BUsecase.Fetch(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(books)
}

func (a *BookHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	ctx := c.Context()
	book, err := a.BUsecase.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(book)
}

func (a *BookHandler) Store(c *fiber.Ctx) error {
	var book domain.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := c.Context()
	err := a.BUsecase.Store(ctx, &book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func (a *BookHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	var book domain.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	book.ID = id

	ctx := c.Context()
	err = a.BUsecase.Update(ctx, &book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(book)
}

func (a *BookHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	ctx := c.Context()
	err = a.BUsecase.Delete(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
