package pgsql

import (
	"log/slog"

	"github.com/cardil/kleio/pkg/clusterlogging"
	"github.com/cardil/kleio/pkg/storage"
	"github.com/cardil/kleio/pkg/storage/inmem"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ storage.Storage = &Storage{}

func (s *Storage) Store(msg *clusterlogging.Message) error {
	return s.pool.AcquireFunc(s.ctx, func(conn *pgxpool.Conn) error {
		sql := `INSERT INTO logs
    (ts, msg, container, image, namespace, pod) 
VALUES
    ($1, $2, $3, $4, $5, $6);`
		_, err := conn.Exec(s.ctx, sql,
			msg.Timestamp.UTC(), msg.Message,
			msg.ContainerName, msg.ContainerImage,
			msg.NamespaceName, msg.PodName,
		)
		return err
	})

}

func (s *Storage) Stats() storage.Stats {
	sts, err := s.stats()
	if err != nil {
		slog.Error("Stats query error", "error", err)
		return nil
	}
	return sts
}

func (s *Storage) stats() (sts storage.Stats, err error) {
	sts = make([]storage.ContainerStat, 0)
	err = s.pool.AcquireFunc(s.ctx, func(conn *pgxpool.Conn) error {
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
		rows, err := conn.Query(s.ctx, sql)
		if rows != nil {
			defer rows.Close()
		}

		if err != nil {
			return err
		}
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
				return err
			}
			sts = append(sts, st)
		}

		if err = rows.Err(); err != nil {
			return err
		}
		return nil
	})

	return
}

func (s *Storage) Download() storage.Artifacts {
	arts, err := s.download()
	if err != nil {
		slog.Error("Download query error", "error", err)
		return nil
	}
	return arts
}

func (s *Storage) download() (arts storage.Artifacts, err error) {
	err = s.pool.AcquireFunc(s.ctx, func(conn *pgxpool.Conn) error {
		sql := `SELECT
    container,
    image,
    namespace,
    pod,
    msg,
    ts
FROM logs ORDER BY id;`
		rows, err := conn.Query(s.ctx, sql)
		if rows != nil {
			defer rows.Close()
		}
		if err != nil {
			return err
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
				return err
			}
			if err = store.Store(e); err != nil {
				return err
			}
		}

		arts = store.Download()
		return nil
	})

	return
}
