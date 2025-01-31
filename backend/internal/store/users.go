package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username , password , email) VALUES($1 , $2 , $3) RETURNING id,created_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)

	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {

	var user User

	query := `Select id,username,email,created_at from users where id= $1`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		return &User{}, ErrNotFound
	}

	return &user, nil

}
