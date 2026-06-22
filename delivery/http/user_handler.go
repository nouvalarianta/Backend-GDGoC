package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nouval/library-app/domain"
)

// UserHandler represent the httphandler for user
type UserHandler struct {
	UUsecase domain.UserUsecase
}

// NewUserHandler will initialize the users/ resources endpoint
func NewUserHandler(app *fiber.App, us domain.UserUsecase) {
	handler := &UserHandler{
		UUsecase: us,
	}
	app.Post("/login", handler.Login)
	app.Post("/users", handler.Store)
}

// Store will store the user by given request body
func (a *UserHandler) Store(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	user := domain.User{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
	}

	ctx := c.Context()
	err := a.UUsecase.Store(ctx, &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login will authenticate the user
func (a *UserHandler) Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username and password are required"})
	}

	ctx := c.Context()
	token, err := a.UUsecase.Login(ctx, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
