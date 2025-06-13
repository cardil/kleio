package collector

import (
	"github.com/cardil/qe-clusterlogging/pkg/storage"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	"v.io/x/lib/vlog"

	"github.com/cardil/qe-clusterlogging/pkg/clusterlogging"
	"gopkg.in/mcuadros/go-syslog.v2"
)

type Collector struct {
	Storage
}

type Storage interface {
	Store(msg *clusterlogging.Message) error
	Stats() storage.Stats
	Download() storage.Artifacts
}

func (c *Collector) Collect(channel syslog.LogPartsChannel) {
	for logParts := range channel {
		if err := c.processLog(logParts); err != nil {
			vlog.Error(err)
			continue
		}
	}
}

func (c *Collector) processLog(logParts format.LogParts) error {
	data, err := clusterlogging.Parse(logParts)
	if err != nil {
		return err
	}
	if err = c.Storage.Store(data); err != nil {
		return err
	}
	return nil
}
