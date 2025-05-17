package testutil

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func InitTestDb(ctx context.Context, dbUrl string) *pgxpool.Pool {
	once.Do(func() {
		var err error // IMPORTANT to avoid variable shadowing on `pool` global variable

		pool, err = pgxpool.New(ctx, dbUrl)
		if err != nil {
			log.Fatalf("failed to create pgx pool: %v", err)
		}
		if err := pool.Ping(ctx); err != nil {
			log.Fatalf("failed to ping pgx pool: %v", err)
		}
	})
	return pool
}

func GetTestDb() *pgxpool.Pool {
	return pool
}
