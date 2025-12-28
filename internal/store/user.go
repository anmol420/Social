package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		u.Username,
		u.Email,
		u.Password,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
