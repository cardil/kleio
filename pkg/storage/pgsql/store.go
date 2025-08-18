package pgsql

import (
	"log/slog"

	"github.com/cardil/kleio/pkg/clusterlogging"
	"github.com/cardil/kleio/pkg/storage"
	"github.com/cardil/kleio/pkg/storage/inmem"
)

var _ storage.Storage = &Storage{}

func (s *Storage) Store(msg *clusterlogging.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sql := `INSERT INTO logs
    (ts, msg, container, image, namespace, pod) 
VALUES
    ($1, $2, $3, $4, $5, $6);`
	_, err := s.conn.Exec(s.ctx, sql,
		msg.Timestamp.UTC(), msg.Message,
		msg.ContainerName, msg.ContainerImage,
		msg.NamespaceName, msg.PodName,
	)
	return err
}

func (s *Storage) Stats() storage.Stats {
	sts, err := s.stats()
	if err != nil {
		slog.Error("Stats query error", "error", err)
		return nil
	}
	return sts
}

func (s *Storage) stats() (storage.Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// See: https://hakibenita.com/sql-group-by-first-last-value
	sql := `select 
       l.container,
       l.image,
       l.namespace,
       l.pod,
       l.count,
       to_timestamp(l.last[2]::double precision / 1000000) as last_ts
from (select container,
             image,
             namespace,
             pod,
             count(1) as count,
             MAX(ARRAY [id, cast(extract(epoch from ts) * 1000000 as bigint)]) as last
      from logs
      group by container, image, namespace, pod) as l;`
	rows, err := s.conn.Query(s.ctx, sql)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		return nil, err
	}
	sts := make([]storage.ContainerStat, 0)
	for rows.Next() {
		st := storage.ContainerStat{}
		if err = rows.Scan(
			&st.ContainerName,
			&st.ContainerImage,
			&st.NamespaceName,
			&st.PodName,
			&st.MessageCount,
			&st.LastMessage,
		); err != nil {
			return nil, err
		}
		sts = append(sts, st)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return sts, nil
}

func (s *Storage) Download() storage.Artifacts {
	arts, err := s.download()
	if err != nil {
		slog.Error("Download query error", "error", err)
		return nil
	}
	return arts
}

func (s *Storage) download() (storage.Artifacts, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sql := `SELECT
    container,
    image,
    namespace,
    pod,
    msg,
    ts
FROM logs ORDER BY id;`
	rows, err := s.conn.Query(s.ctx, sql)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}

	store := inmem.NewStore()
	for rows.Next() {
		e := &clusterlogging.Message{}
		if err = rows.Scan(
			&e.ContainerName,
			&e.ContainerImage,
			&e.NamespaceName,
			&e.PodName,
			&e.Message,
			&e.Timestamp,
		); err != nil {
			return nil, err
		}
		if err = store.Store(e); err != nil {
			return nil, err
		}
	}

	return store.Download(), nil
}
