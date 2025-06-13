package main

import (
	"github.com/cardil/qe-clusterlogging/pkg/collector"
	"github.com/cardil/qe-clusterlogging/pkg/storage/inmem"
	"github.com/cardil/qe-clusterlogging/pkg/syslog"
	"v.io/x/lib/vlog"
)

func main() {
	c := &collector.Collector{Storage: inmem.NewStore()}
	waiter, err := syslog.Serve(c.Collect)
	if err != nil {
		vlog.Fatal(err)
	}
	waiter.Wait()
}
