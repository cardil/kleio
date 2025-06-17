package collector_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/cardil/kleio/internal/collector"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
)

func TestServe(t *testing.T) {
	server := collector.Serve()
	go func() {
		require.NoError(t, server.Run())
	}()
	defer func() {
		require.NoError(t, server.Kill())
	}()
	require.NoError(t, waitForPortOpen(8514))
	require.NoError(t, waitForPortOpen(8080))
}

func TestServeOrDie(t *testing.T) {
	t.Setenv("PORT", "70000")
	var gotRetcode *int
	collector.ServeOrDie(func(rc int) {
		gotRetcode = &rc
	})
	require.NotNil(t, gotRetcode)
	require.Equal(t, 135, *gotRetcode)
}

func waitForPortOpen(port int) error {
	interval := time.Millisecond
	timeout := 5 * time.Second
	return wait.PollUntilContextTimeout(context.TODO(), interval, timeout, true, func(ctx context.Context) (done bool, err error) {
		done = isPortOpen(port)
		return
	})
}

func isPortOpen(port int) bool {
	timeout := time.Second
	addr := net.JoinHostPort("127.0.0.1", fmt.Sprint(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer func(c io.Closer) {
			_ = c.Close()
		}(conn)
	}
	return conn != nil
}
