package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// SystemCheckRepository is a contract interface for system checks.
type SystemCheckRepository interface {
	CheckPostgresql(ctx context.Context) error
	CheckRedis(ctx context.Context, rdb *redis.Client) error
}

// systemCheckRepository implements SystemCheckRepository interface
type systemCheckRepository struct {
	db *pgxpool.Pool
}

func NewSystemCheckRepository(db *pgxpool.Pool) SystemCheckRepository {
	return &systemCheckRepository{
		db: db,
	}
}

func (r *systemCheckRepository) CheckPostgresql(ctx context.Context) error {
	if err := r.db.Ping(ctx); err != nil {
		return err
	}
	return nil
}

func (r *systemCheckRepository) CheckRedis(ctx context.Context, rdb *redis.Client) error {
	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}
