package collector_test

import (
	"testing"

	"github.com/cardil/kleio/pkg/collector"
	"github.com/cardil/kleio/pkg/storage/inmem"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestCollector_Collect(t *testing.T) {
	t.Parallel()
	st := inmem.NewStore()
	coll := collector.Collector{Storage: st}
	ch := make(chan format.LogParts)
	go coll.Collect(ch)
	ch <- map[string]interface{}{
		"message": `{"message": "foo", "timestamp":"2025-06-17T18:52:12Z"}`,
	}
	close(ch)
	stts := st.Stats()
	assert.Len(t, stts, 1)
	first := stts[0]
	assert.Equal(t, 1, first.MessageCount)
}
