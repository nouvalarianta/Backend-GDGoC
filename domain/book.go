package domain

import (
	"context"
	"time"
)

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookRepository interface {
	Fetch(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id int) (Book, error)
	Store(ctx context.Context, b *Book) error
	Update(ctx context.Context, b *Book) error
	Delete(ctx context.Context, id int) error
	UpdateAvailability(ctx context.Context, id int, isAvailable bool) error
}

type BookUsecase interface {
	Fetch(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id int) (Book, error)
	Store(ctx context.Context, b *Book) error
	Update(ctx context.Context, b *Book) error
	Delete(ctx context.Context, id int) error
}
