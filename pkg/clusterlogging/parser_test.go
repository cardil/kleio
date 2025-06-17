package clusterlogging_test

import (
	"testing"
	"time"

	"github.com/cardil/qe-clusterlogging/pkg/clusterlogging"
	"github.com/cardil/qe-clusterlogging/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mcuadros/go-syslog.v2/format"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	dt := time.Date(2025, 6, 17, 9, 52, 57,
		46_000_000, time.FixedZone("", 0))
	tcs := []parseTestCase{{
		name:    "empty",
		wantErr: clusterlogging.ErrInvalidFormat,
	}, {
		name:    "log with []byte message",
		wantErr: clusterlogging.ErrInvalidFormat,
		input: map[string]interface{}{
			"message": make([]byte, 10),
		},
	}, {
		name:    "log with bool as message",
		wantErr: clusterlogging.ErrInvalidFormat,
		input: map[string]interface{}{
			"message": true,
		},
	}, {
		name:    "log with string message",
		wantErr: clusterlogging.ErrInvalidFormat,
		input: map[string]interface{}{
			"message": "Alice in Wonderland",
		},
	}, {
		name: "legal message",
		want: &clusterlogging.Message{
			Timestamp: dt,
			Message:   "INFO msg",
			ContainerInfo: kubernetes.ContainerInfo{
				ContainerName:  "user",
				ContainerImage: "example.org/foo",
				NamespaceName:  "default",
				PodName:        "foo",
			},
		},
		input: map[string]interface{}{
			"message": `{
	"message": "INFO msg",
	"timestamp": "2025-06-17T09:52:57.046+00:00",
	"kubernetes": {
		"pod_name": "foo",
		"namespace_name": "default",
		"container_name": "user",
		"container_image": "example.org/foo"
	}
}`,
		},
	}}
	for _, tc := range tcs {
		t.Run(tc.name, tc.test)
	}
}

type parseTestCase struct {
	name    string
	input   format.LogParts
	want    *clusterlogging.Message
	wantErr error
}

func (tc parseTestCase) test(t *testing.T) {
	msg, err := clusterlogging.Parse(tc.input)
	require.ErrorIs(t, err, tc.wantErr)
	assert.EqualValues(t, tc.want, msg)
}
