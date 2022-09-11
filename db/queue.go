package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	etcd "go.etcd.io/etcd/client/v3"
	concurrency "go.etcd.io/etcd/client/v3/concurrency"
)

type Queue struct {
	Pointer int `json:"pointer"`
	Size    int `json:"size"`
	MaxSize int `json:"max_size"`
}

var ErrEmptyDBQueue error = errors.New("empty db queue")

func (db *DB) Dequeue(v interface{}, key string) error {
	_, ok := db.qLocks[key]
	if !ok {
		db.qLocks[key] = &sync.Mutex{}
	}
	db.qLocks[key].Lock()
	defer db.qLocks[key].Unlock()

	if s, err := concurrency.NewSession(db.client); err != nil {
		return err
	} else {
		defer s.Close()
		l := concurrency.NewLocker(s, key)
		l.Lock()
		defer l.Unlock()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		q := &Queue{}

		if r, err := s.Client().Get(ctx, key, etcd.WithPrefix()); err != nil {
			return err
		} else if r.Count == 0 {
			return ErrEmptyDBQueue
		} else if err := json.Unmarshal(r.Kvs[0].Value, q); err != nil {
			return err
		} else if r.Count == 1 {
			return ErrEmptyDBQueue
		} else if err := json.Unmarshal(r.Kvs[q.Pointer+1].Value, v); err != nil {
			return err
		} else {
			iKey := fmt.Sprintf("%s/%d", key, q.Pointer)
			s.Client().Delete(ctx, iKey)

			q.Pointer++
			if q.Pointer == q.Size {
				q.Pointer = 0
				q.Size = 0
			}
			d, _ := json.Marshal(q)
			if _, err := s.Client().Put(ctx, key, string(d)); err != nil {
				return err
			}
			return nil
		}
	}
}

var ErrDBQueueOverflow error = errors.New("db queue is overflow")

func (db *DB) Enqueue(v interface{}, key string) error {
	_, ok := db.qLocks[key]
	if !ok {
		db.qLocks[key] = &sync.Mutex{}
	}
	db.qLocks[key].Lock()
	defer db.qLocks[key].Unlock()

	if s, err := concurrency.NewSession(db.client); err != nil {
		return err
	} else {
		defer s.Close()
		l := concurrency.NewLocker(s, key)
		l.Lock()
		defer l.Unlock()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		q := &Queue{}

		if r, err := s.Client().Get(ctx, key, etcd.WithPrefix()); err != nil {
			return err
		} else if err := json.Unmarshal(r.Kvs[0].Value, q); err != nil {
			return err
		} else {
			q.Size++
			if q.Size > q.MaxSize {
				return ErrDBQueueOverflow
			}

			iKey := fmt.Sprintf("%s/%d", key, q.Pointer)
			dq, _ := json.Marshal(q)
			dv, _ := json.Marshal(v)
			if _, err := s.Client().Put(ctx, iKey, string(dv)); err != nil {
				return err
			} else if _, err := s.Client().Put(ctx, key, string(dq)); err != nil {
				return err
			}
			return nil
		}
	}
}
