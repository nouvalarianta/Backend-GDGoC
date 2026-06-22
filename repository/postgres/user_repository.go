package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nouval/library-app/domain"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

// NewPostgresUserRepository will create an object that represent the user.Repository interface
func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{Conn: conn}
}

func (m *mysqlUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.User, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.User
	for rows.Next() {
		t := domain.User{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Password,
			&t.Name,
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

func (m *mysqlUserRepository) Fetch(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id, username, password, name, created_at, updated_at FROM users`
	return m.fetch(ctx, query)
}

func (m *mysqlUserRepository) GetByID(ctx context.Context, id int) (domain.User, error) {
	query := `SELECT id, username, password, name, created_at, updated_at FROM users WHERE id = $1`
	res, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}
	if len(res) == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return res[0], nil
}

func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `SELECT id, username, password, name, created_at, updated_at FROM users WHERE username = $1`
	res, err := m.fetch(ctx, query, username)
	if err != nil {
		return domain.User{}, err
	}
	if len(res) == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return res[0], nil
}

func (m *mysqlUserRepository) Store(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (username, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := m.Conn.QueryRowContext(ctx, query, u.Username, u.Password, u.Name, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	return err
}

func (m *mysqlUserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE users SET name = $1, password = $2, updated_at = $3 WHERE id = $4`
	_, err := m.Conn.ExecContext(ctx, query, u.Name, u.Password, u.UpdatedAt, u.ID)
	return err
}

func (m *mysqlUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := m.Conn.ExecContext(ctx, query, id)
	return err
}
