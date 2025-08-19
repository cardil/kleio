package facade_test

import (
	"context"
	"testing"

	"github.com/cardil/kleio/pkg/storage/facade"
	"github.com/cardil/kleio/pkg/storage/pgsql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreFacade(t *testing.T) {
	ctx := context.TODO()
	t.Setenv("DATABASE_URI", "")
	store, err := facade.NewStore(ctx)
	require.NoError(t, err)
	assert.NotNil(t, store)
}

func TestStorePostgresFacade(t *testing.T) {
	ctx := context.TODO()
	t.Setenv("DATABASE_URI", "postgresql://myuser:bad@127.0.0.1:15432/bad")
	_, err := facade.NewStore(ctx)
	require.ErrorIs(t, err, pgsql.ErrCantConnect)
}
