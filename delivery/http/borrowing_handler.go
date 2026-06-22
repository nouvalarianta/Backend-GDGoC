package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nouval/library-app/delivery/http/middleware"
	"github.com/nouval/library-app/domain"
)

// BorrowingHandler represent the httphandler for borrowing
type BorrowingHandler struct {
	BUsecase domain.BorrowingUsecase
}

// NewBorrowingHandler will initialize the borrowings/ resources endpoint
func NewBorrowingHandler(app *fiber.App, us domain.BorrowingUsecase, jwtSecret string) {
	handler := &BorrowingHandler{
		BUsecase: us,
	}

	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	protected := app.Group("/borrowings", authMiddleware)

	protected.Post("/borrow", handler.BorrowBook)
	protected.Post("/return/:id", handler.ReturnBook)
	protected.Get("/user", handler.GetUserBorrowings)
	protected.Get("/", handler.FetchAll)
}

func (a *BorrowingHandler) BorrowBook(c *fiber.Ctx) error {
	type BorrowRequest struct {
		BookID int `json:"book_id"`
	}

	var req BorrowRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(int)
	ctx := c.Context()

	err := a.BUsecase.BorrowBook(ctx, userID, req.BookID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "book borrowed successfully"})
}

func (a *BorrowingHandler) ReturnBook(c *fiber.Ctx) error {
	borrowingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid borrowing ID"})
	}

	userID := c.Locals("user_id").(int)
	ctx := c.Context()

	err = a.BUsecase.ReturnBook(ctx, borrowingID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "book returned successfully"})
}

func (a *BorrowingHandler) GetUserBorrowings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	ctx := c.Context()

	borrowings, err := a.BUsecase.GetUserBorrowings(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(borrowings)
}

func (a *BorrowingHandler) FetchAll(c *fiber.Ctx) error {
	ctx := c.Context()

	borrowings, err := a.BUsecase.FetchAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(borrowings)
}
