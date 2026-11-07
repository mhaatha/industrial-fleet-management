package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// SystemCheckRepository is a contract interface for system checks.
type SystemCheckRepository interface {
	CheckPostgresql(ctx context.Context) error
	CheckRedis(ctx context.Context) error
}

// systemCheckRepository implements SystemCheckRepository interface
type systemCheckRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewSystemCheckRepository(db *pgxpool.Pool, rdb *redis.Client) SystemCheckRepository {
	return &systemCheckRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *systemCheckRepository) CheckPostgresql(ctx context.Context) error {
	if err := r.db.Ping(ctx); err != nil {
		return err
	}
	return nil
}

func (r *systemCheckRepository) CheckRedis(ctx context.Context) error {
	if err := r.rdb.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}
