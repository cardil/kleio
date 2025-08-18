package pgsql_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cardil/kleio/pkg/clusterlogging"
	"github.com/cardil/kleio/pkg/kubernetes"
	"github.com/cardil/kleio/pkg/storage/pgsql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testpgsql "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPgsqlStorage(t *testing.T) {
	ctx := t.Context()

	// TODO: Remove if fixed
	//       See: https://github.com/testcontainers/testcontainers-go/issues/3262
	if strings.Contains(os.Getenv("DOCKER_HOST"), "podman.sock") {
		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	}
	pgsqlContainer, err := testpgsql.Run(ctx,
		"docker.io/library/postgres:17-alpine",
		testpgsql.WithDatabase("mydb"),
		testpgsql.WithUsername("myuser"),
		testpgsql.WithPassword("mypassword"),
		testpgsql.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pgsqlContainer.Terminate(ctx))
	}()
	uri := pgsqlContainer.MustConnectionString(ctx)
	cfg, err := pgxpool.ParseConfig(uri)
	require.NoError(t, err)
	store, err := pgsql.NewStore(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, store.Close())
	}()

	msg := &clusterlogging.Message{
		Timestamp: time.Now(),
		Message:   "test message",
		ContainerInfo: kubernetes.ContainerInfo{
			ContainerName:  "user",
			ContainerImage: "example.org/test",
			NamespaceName:  "default",
			PodName:        "acme",
		},
	}
	err = store.Store(msg)
	require.NoError(t, err)

	sts := store.Stats()
	assert.NotNil(t, sts)
	assert.Len(t, sts, 1)
	first := sts[0]
	assert.Equal(t, msg.FullName(), first.FullName())
	assert.Equal(t, 1, first.MessageCount)
	assert.InDelta(t, msg.Timestamp.UnixMilli(), first.LastMessage.UnixMilli(), 10)
}
