package pgsql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCantConnect = errors.New("can't connect to database")

func NewStore(ctx context.Context, cfg *pgxpool.Config) (*Storage, error) {
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantConnect, err)
	}
	st := &Storage{
		pool: pool,
		ctx:  ctx,
	}

	err = st.init()
	connString := cfg.ConnString()
	if u, err := url.Parse(connString); err != nil {
		connString = "xxxxx"
	} else {
		connString = u.Redacted()
	}
	slog.Info("Using PostgreSQL store", "conn", connString)
	return st, err
}

func NewStoreFromEnvironment(ctx context.Context) (*Storage, error) {
	uri := os.Getenv("DATABASE_URI")
	uri = strings.Replace(uri, "postgres://", "postgresql://", 1)
	if !strings.HasPrefix(uri, "postgresql://") {
		return nil, nil
	}
	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantConnect, err)
	}
	return NewStore(ctx, cfg)
}

type Storage struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func (s *Storage) Close() error {
	s.pool.Close()
	return nil
}

func (s *Storage) init() error {
	if err := s.pool.AcquireFunc(s.ctx, func(c *pgxpool.Conn) error {
		return s.initWithConn(c)
	}); err != nil {
		return err
	}
	return nil
}

func (s *Storage) initWithConn(conn *pgxpool.Conn) (err error) {
	_, err = conn.Exec(s.ctx, `CREATE TABLE IF NOT EXISTS logs (
  id        integer primary key generated always as identity,
  ts        timestamp not null,
  msg       varchar(32768) not null,
  container varchar(80) not null,
  image     varchar(500) not null,
  namespace varchar(80) not null,
  pod       varchar(80) not null
);`)
	_, perr := conn.Exec(s.ctx, `CREATE INDEX IF NOT EXISTS 
container_info_index ON logs (container, image, namespace, pod, id DESC, ts DESC)`)
	err = errors.Join(err, perr)
	return
}
