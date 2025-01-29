package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/solrac97gr/cqrs/models"
	"github.com/solrac97gr/cqrs/repository"
)

type PostgresRepository struct {
	db *sql.DB
}

var _ repository.Repository = &PostgresRepository{}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO feeds (id, title, description) VALUES ($1, $2, $3)", feed.ID, feed.Title, feed.Description)
	return err
}

func (r *PostgresRepository) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, description, created_at FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []*models.Feed{}
	for rows.Next() {
		var feed models.Feed
		if err := rows.Scan(&feed.ID, &feed.Title, &feed.Description, &feed.CreatedAt); err != nil {
			return nil, err
		}
		feeds = append(feeds, &feed)
	}
	return feeds, nil
}
