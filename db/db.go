package db

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/spf13/viper"
	etcd "go.etcd.io/etcd/client/v3"
	concurrency "go.etcd.io/etcd/client/v3/concurrency"
)

type DB struct {
	client *etcd.Client
	qLocks map[string]sync.Locker
}

var ErrValueNotFound = errors.New("value not found")

func (db *DB) Get(ctx context.Context, key string, v interface{}) error {
	if r, err := db.client.Get(ctx, key); err != nil {
		return err
	} else if len(r.Kvs) == 0 {
		return ErrValueNotFound
	} else {
		return json.Unmarshal(r.Kvs[0].Value, v)
	}
}

func (db *DB) GetList(ctx context.Context, prefix string, v interface{}, f func(interface{}) bool) error {
	if r, err := db.client.Get(ctx, prefix, etcd.WithPrefix()); err != nil {
		return nil
	} else {
		for _, kv := range r.Kvs {
			json.Unmarshal(kv.Value, v)
			if next := f(v); !next {
				break
			}
		}
		return nil
	}
}

func (db *DB) GetWithPrefix(ctx context.Context, prefix string) []interface{} {
	if r, err := db.client.Get(ctx, prefix, etcd.WithPrefix()); err != nil {
		return nil
	} else {
		l := make([]interface{}, 0)
		for _, kv := range r.Kvs {
			l = append(l, kv.Value)
		}
		return l
	}
}

func (db *DB) GetString(ctx context.Context, key string) (string, error) {
	if r, err := db.client.Get(ctx, key); err != nil {
		return "", err
	} else {
		return string(r.Kvs[0].Value), nil
	}
}

func (db *DB) Lock(ctx context.Context, key string, f func()) error {
	if s, err := concurrency.NewSession(db.client); err != nil {
		return err
	} else {
		defer s.Close()
		l := concurrency.NewLocker(s, key)
		l.Lock()
		defer l.Unlock()
		f()
		return nil
	}
}

func (db *DB) Put(ctx context.Context, key string, v interface{}) error {
	if d, err := json.Marshal(v); err != nil {
		return err
	} else {
		_, err := db.client.Put(ctx, key, string(d))
		return err
	}
}

func NewDB() (*DB, error) {
	endpoints := viper.GetStringSlice("db.endpoints")
	client, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return &DB{
		client: client,
	}, nil
}
