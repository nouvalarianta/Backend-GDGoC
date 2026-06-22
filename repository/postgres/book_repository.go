package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nouval/library-app/domain"
)

type postgresBookRepository struct {
	Conn *sql.DB
}

func NewPostgresBookRepository(conn *sql.DB) domain.BookRepository {
	return &postgresBookRepository{Conn: conn}
}

func (m *postgresBookRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Book, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Book
	for rows.Next() {
		t := domain.Book{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Author,
			&t.ISBN,
			&t.IsAvailable,
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

func (m *postgresBookRepository) Fetch(ctx context.Context) ([]domain.Book, error) {
	query := `SELECT id, title, author, isbn, is_available, created_at, updated_at FROM books`
	return m.fetch(ctx, query)
}

func (m *postgresBookRepository) GetByID(ctx context.Context, id int) (domain.Book, error) {
	query := `SELECT id, title, author, isbn, is_available, created_at, updated_at FROM books WHERE id = $1`
	res, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Book{}, err
	}
	if len(res) == 0 {
		return domain.Book{}, errors.New("book not found")
	}
	return res[0], nil
}

func (m *postgresBookRepository) Store(ctx context.Context, b *domain.Book) error {
	query := `INSERT INTO books (title, author, isbn, is_available, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := m.Conn.QueryRowContext(ctx, query, b.Title, b.Author, b.ISBN, b.IsAvailable, b.CreatedAt, b.UpdatedAt).Scan(&b.ID)
	return err
}

func (m *postgresBookRepository) Update(ctx context.Context, b *domain.Book) error {
	query := `UPDATE books SET title = $1, author = $2, isbn = $3, is_available = $4, updated_at = $5 WHERE id = $6`
	_, err := m.Conn.ExecContext(ctx, query, b.Title, b.Author, b.ISBN, b.IsAvailable, b.UpdatedAt, b.ID)
	return err
}

func (m *postgresBookRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := m.Conn.ExecContext(ctx, query, id)
	return err
}

func (m *postgresBookRepository) UpdateAvailability(ctx context.Context, id int, isAvailable bool) error {
	query := `UPDATE books SET is_available = $1 WHERE id = $2`
	_, err := m.Conn.ExecContext(ctx, query, isAvailable, id)
	return err
}
