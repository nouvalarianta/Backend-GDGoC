package domain

import (
	"context"
	"time"
)

const (
	BorrowStatusBorrowed = "borrowed"
	BorrowStatusReturned = "returned"
)

type Borrowing struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	BookID     int        `json:"book_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	ReturnDate *time.Time `json:"return_date"` // Can be null if not returned yet
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	User *User `json:"user,omitempty"`
	Book *Book `json:"book,omitempty"`
}

type BorrowingRepository interface {
	Fetch(ctx context.Context) ([]Borrowing, error)
	GetByID(ctx context.Context, id int) (Borrowing, error)
	GetByUserID(ctx context.Context, userID int) ([]Borrowing, error)
	Store(ctx context.Context, b *Borrowing) error
	Update(ctx context.Context, b *Borrowing) error
}

type BorrowingUsecase interface {
	BorrowBook(ctx context.Context, userID int, bookID int) error
	ReturnBook(ctx context.Context, borrowingID int, userID int) error
	GetUserBorrowings(ctx context.Context, userID int) ([]Borrowing, error)
	FetchAll(ctx context.Context) ([]Borrowing, error)
}
