package db

import (
	"context"
	"fmt"
	"net"
	"sync"

	"signal/sip"

	"github.com/google/uuid"
	concurrency "go.etcd.io/etcd/client/v3/concurrency"
)

type Subscriber struct {
	db      *DB
	session *concurrency.Session
	L       sync.Locker
	ID      uuid.UUID
	URI     sip.URI  `json:"address"`
	IsBusy  bool     `json:"is_busy"`
	Addr    net.Addr `json:"addr"`
}

func (db *DB) CaptureSubscriber(ctx context.Context, id uuid.UUID) (*Subscriber, error) {
	if session, err := concurrency.NewSession(db.client); err != nil {
		return nil, err
	} else {
		s := &Subscriber{
			db:      db,
			ID:      id,
			session: session,
		}
		key := fmt.Sprintf("/subscriber/%s", id.String())
		s.L = concurrency.NewLocker(session, key)
		s.L.Lock()

		db.Get(ctx, key, s)

		return s, nil
	}
}

func (s *Subscriber) Save(ctx context.Context) error {
	key := fmt.Sprintf("/subscriber/%s", s.ID.String())
	return s.db.Put(ctx, key, s)
}
