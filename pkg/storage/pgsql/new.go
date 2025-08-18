package pgsql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
)

var ErrCantConnect = errors.New("can't connect to database")

func NewStore(ctx context.Context, cfg *pgx.ConnConfig) (*Storage, error) {
	connString := cfg.ConnString()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantConnect, err)
	}
	st := &Storage{
		conn: conn,
		ctx:  ctx,
	}
	connString = strings.Replace(connString, cfg.Password, "***", 1)
	err = st.init()
	slog.Info("Using PostgreSQL store", "conn", connString)
	return st, err
}

func NewStoreFromEnvironment(ctx context.Context) (*Storage, error) {
	uri := os.Getenv("DATABASE_URI")
	uri = strings.Replace(uri, "postgres://", "postgresql://", 1)
	if !strings.HasPrefix(uri, "postgresql://") {
		return nil, nil
	}
	cfg, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantConnect, err)
	}
	return NewStore(ctx, cfg)
}

type Storage struct {
	conn *pgx.Conn
	ctx  context.Context
	mu   sync.RWMutex
}

func (s *Storage) Close() error {
	return s.conn.Close(context.Background())
}

func (s *Storage) init() (err error) {
	_, err = s.conn.Exec(s.ctx, `CREATE TABLE IF NOT EXISTS logs (
  id        integer primary key generated always as identity,
  ts        timestamp not null,
  msg       varchar(4096) not null,
  container varchar(80) not null,
  image     varchar(500) not null,
  namespace varchar(80) not null,
  pod       varchar(80) not null
);`)
	_, perr := s.conn.Exec(s.ctx, `CREATE INDEX IF NOT EXISTS 
container_info_index ON logs (container, image, namespace, pod)`)
	err = errors.Join(err, perr)
	return
}
