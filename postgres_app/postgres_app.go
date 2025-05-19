package postgres_app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string // e.g., "disable", "require"
}

type PostgresApp struct {
	cfg *Config
	DB  *pgxpool.Pool
}

func New(cfg *Config) *PostgresApp {
	return &PostgresApp{cfg: cfg}
}

func (p *PostgresApp) Start(ctx context.Context) error {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=10",
		p.cfg.Username,
		p.cfg.Password,
		p.cfg.Host,
		p.cfg.Port,
		p.cfg.DBName,
		p.cfg.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	p.DB = pool
	return nil
}

func (p *PostgresApp) Stop(ctx context.Context) error {
	if p.DB != nil {
		p.DB.Close()
	}
	return nil
}
