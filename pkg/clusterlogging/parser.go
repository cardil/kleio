package clusterlogging

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var ErrInvalidFormat = errors.New("invalid format")

func Parse(logParts format.LogParts) (*Message, error) {
	raw, ok := logParts["message"]
	if !ok {
		return nil, fmt.Errorf("%w: message field missing",
			ErrInvalidFormat)
	}
	var bytes []byte
	switch v := raw.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return nil, fmt.Errorf("%w: unexpected message type %T",
			ErrInvalidFormat, v)
	}
	data := &Message{}
	if err := json.Unmarshal(bytes, data); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	return data, nil
}
