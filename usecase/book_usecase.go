package usecase

import (
	"context"
	"time"

	"github.com/nouval/library-app/domain"
)

type bookUsecase struct {
	bookRepo       domain.BookRepository
	contextTimeout time.Duration
}

// NewBookUsecase will create new an bookUsecase object representation of domain.BookUsecase interface
func NewBookUsecase(b domain.BookRepository, timeout time.Duration) domain.BookUsecase {
	return &bookUsecase{
		bookRepo:       b,
		contextTimeout: timeout,
	}
}

func (u *bookUsecase) Fetch(c context.Context) ([]domain.Book, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.bookRepo.Fetch(ctx)
}

func (u *bookUsecase) GetByID(c context.Context, id int) (domain.Book, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.bookRepo.GetByID(ctx, id)
}

func (u *bookUsecase) Store(c context.Context, m *domain.Book) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.IsAvailable = true

	return u.bookRepo.Store(ctx, m)
}

func (u *bookUsecase) Update(c context.Context, m *domain.Book) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	m.UpdatedAt = time.Now()
	return u.bookRepo.Update(ctx, m)
}

func (u *bookUsecase) Delete(c context.Context, id int) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.bookRepo.Delete(ctx, id)
}
