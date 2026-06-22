package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/nouval/library-app/domain"
)

type borrowingUsecase struct {
	borrowingRepo  domain.BorrowingRepository
	bookRepo       domain.BookRepository
	contextTimeout time.Duration
}

// NewBorrowingUsecase will create new an borrowingUsecase object representation of domain.BorrowingUsecase interface
func NewBorrowingUsecase(br domain.BorrowingRepository, b domain.BookRepository, timeout time.Duration) domain.BorrowingUsecase {
	return &borrowingUsecase{
		borrowingRepo:  br,
		bookRepo:       b,
		contextTimeout: timeout,
	}
}

func (u *borrowingUsecase) BorrowBook(c context.Context, userID int, bookID int) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// 1. cek kesediaan buku
	book, err := u.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}

	if !book.IsAvailable {
		return errors.New("book is not available for borrowing")
	}

	// 2. create data peminjaman
	borrowing := &domain.Borrowing{
		UserID:     userID,
		BookID:     bookID,
		BorrowDate: time.Now(),
		Status:     domain.BorrowStatusBorrowed,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = u.borrowingRepo.Store(ctx, borrowing)
	if err != nil {
		return err
	}

	// 3. Update status buku
	return u.bookRepo.UpdateAvailability(ctx, bookID, false)
}

func (u *borrowingUsecase) ReturnBook(c context.Context, borrowingID int, userID int) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// 1. Fetch data peminjaman
	borrowing, err := u.borrowingRepo.GetByID(ctx, borrowingID)
	if err != nil {
		return err
	}

	if borrowing.UserID != userID {
		return errors.New("unauthorized to return this book")
	}

	if borrowing.Status == domain.BorrowStatusReturned {
		return errors.New("book is already returned")
	}

	// 2. Update data peminjaman
	now := time.Now()
	borrowing.ReturnDate = &now
	borrowing.Status = domain.BorrowStatusReturned
	borrowing.UpdatedAt = now

	err = u.borrowingRepo.Update(ctx, &borrowing)
	if err != nil {
		return err
	}

	// 3. Update status buku jadi tersedia
	return u.bookRepo.UpdateAvailability(ctx, borrowing.BookID, true)
}

func (u *borrowingUsecase) GetUserBorrowings(c context.Context, userID int) ([]domain.Borrowing, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.borrowingRepo.GetByUserID(ctx, userID)
}

func (u *borrowingUsecase) FetchAll(c context.Context) ([]domain.Borrowing, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.borrowingRepo.Fetch(ctx)
}
