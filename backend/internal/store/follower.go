package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

const QueryTimeoutDuration = time.Second * 3

type Follower struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt  int64 `json:"created_at"`
}

type FollowStore struct {
	db *sql.DB
}

func (s *FollowStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("Already Followed")
		}
	}

	return nil
}

func (s *FollowStore) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		DELETE FROM followers 
		WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}
