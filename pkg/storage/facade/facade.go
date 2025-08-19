package facade

import (
	"context"

	"github.com/cardil/kleio/pkg/storage"
	"github.com/cardil/kleio/pkg/storage/inmem"
	"github.com/cardil/kleio/pkg/storage/pgsql"
)

func NewStore(ctx context.Context) (storage.Storage, error) {
	pgstore, err := pgsql.NewStoreFromEnvironment(ctx)
	if err != nil {
		return nil, err
	}
	if pgstore != nil {
		return pgstore, nil
	}
	return inmem.NewStore(), nil
}
