package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	etcd "go.etcd.io/etcd/client/v3"
)

type DB interface {
	Get(context.Context, string, interface{}) error
	Put(context.Context, string, interface{}) error
	Delete(context.Context, string) error
}

var ErrValueNotFound = errors.New("value not found")

type SQLiteDB struct {
	db *sql.DB
}

func (driver *SQLiteDB) Get(ctx context.Context, key string, value interface{}) error {
	var rawData []byte
	if stmt, err := driver.db.Prepare("SELECT value FROM store WHERE key=?"); err != nil {
		return err
	} else if err := stmt.QueryRow(key).Scan(rawData); err != nil {
		stmt.Close()
		return err
	} else {
		defer stmt.Close()
		return json.Unmarshal(rawData, value)
	}
}

func (driver *SQLiteDB) Put(ctx context.Context, key string, value interface{}) error {
	if rawData, err := json.Marshal(value); err != nil {
		return err
	} else {
		if tx, err := driver.db.Begin(); err != nil {
			return err
		} else if stmt, err := tx.Prepare("INSERT INTO store(key, value) values(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value"); err != nil {
			return err
		} else {
			defer stmt.Close()
			if _, err := stmt.Exec(key, rawData); err != nil {
				return err
			} else if err := tx.Commit(); err != nil {
				return err
			}
			return nil
		}
	}
}

func (driver *SQLiteDB) Delete(ctx context.Context, key string) error {
	_, err := driver.db.Exec("DELETE FROM store WHERE key=?", key)
	return err
}

func NewSQLiteDB(filepath string) (*SQLiteDB, error) {
	if db, err := sql.Open("sqlite3", filepath); err != nil {
		return nil, err
	} else {
		return &SQLiteDB{
			db: db,
		}, nil
	}
}

type ETCDDB struct {
	client *etcd.Client
}

func (driver *ETCDDB) Get(ctx context.Context, key string, v interface{}) error {
	if r, err := driver.client.Get(ctx, key); err != nil {
		return err
	} else if len(r.Kvs) == 0 {
		return ErrValueNotFound
	} else {
		return json.Unmarshal(r.Kvs[0].Value, v)
	}
}

func (driver *ETCDDB) Put(ctx context.Context, key string, value []byte) error {
	if d, err := json.Marshal(value); err != nil {
		return err
	} else {
		_, err := driver.client.Put(ctx, key, string(d))
		return err
	}
}

func (driver *ETCDDB) Delete(ctx context.Context, key string) error {
	if _, err := driver.client.Delete(ctx, key); err != nil {
		return err
	}
	return nil
}

func NewETCDDB() (*ETCDDB, error) {
	endpoints := viper.GetStringSlice("db.endpoints")
	client, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return &ETCDDB{
		client: client,
	}, nil
}
