package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nouval/library-app/domain"
)

type postgresBorrowingRepository struct {
	Conn *sql.DB
}

// NewPostgresBorrowingRepository will create an object that represent the borrowing.Repository interface
func NewPostgresBorrowingRepository(conn *sql.DB) domain.BorrowingRepository {
	return &postgresBorrowingRepository{Conn: conn}
}

func (m *postgresBorrowingRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Borrowing, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Borrowing
	for rows.Next() {
		t := domain.Borrowing{}
		err = rows.Scan(
			&t.ID,
			&t.UserID,
			&t.BookID,
			&t.BorrowDate,
			&t.ReturnDate,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *postgresBorrowingRepository) Fetch(ctx context.Context) ([]domain.Borrowing, error) {
	query := `SELECT id, user_id, book_id, borrow_date, return_date, status, created_at, updated_at FROM borrowings`
	return m.fetch(ctx, query)
}

func (m *postgresBorrowingRepository) GetByID(ctx context.Context, id int) (domain.Borrowing, error) {
	query := `SELECT id, user_id, book_id, borrow_date, return_date, status, created_at, updated_at FROM borrowings WHERE id = $1`
	res, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Borrowing{}, err
	}
	if len(res) == 0 {
		return domain.Borrowing{}, errors.New("borrowing record not found")
	}
	return res[0], nil
}

func (m *postgresBorrowingRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Borrowing, error) {
	query := `SELECT id, user_id, book_id, borrow_date, return_date, status, created_at, updated_at FROM borrowings WHERE user_id = $1`
	return m.fetch(ctx, query, userID)
}

func (m *postgresBorrowingRepository) Store(ctx context.Context, b *domain.Borrowing) error {
	query := `INSERT INTO borrowings (user_id, book_id, borrow_date, return_date, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := m.Conn.QueryRowContext(ctx, query, b.UserID, b.BookID, b.BorrowDate, b.ReturnDate, b.Status, b.CreatedAt, b.UpdatedAt).Scan(&b.ID)
	return err
}

func (m *postgresBorrowingRepository) Update(ctx context.Context, b *domain.Borrowing) error {
	query := `UPDATE borrowings SET user_id = $1, book_id = $2, borrow_date = $3, return_date = $4, status = $5, updated_at = $6 WHERE id = $7`
	_, err := m.Conn.ExecContext(ctx, query, b.UserID, b.BookID, b.BorrowDate, b.ReturnDate, b.Status, b.UpdatedAt, b.ID)
	return err
}
