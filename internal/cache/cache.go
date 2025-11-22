package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AletisSearch/aletis/internal/db"
	"github.com/jackc/pgx/v5"
)

type CacheMarshal interface {
	UnmarshalMsg(bts []byte) (o []byte, err error)
	MarshalMsg(b []byte) (o []byte, err error)
}

var ErrNotFoundInCache = errors.New("cache result not found")
var ErrOldCache = errors.New("cache result old")

type Cache[T any, PT interface {
	*T
	CacheMarshal
}] struct {
	db *db.Queries
}

func New[T any, PT interface {
	*T
	CacheMarshal
}](db *db.Queries) *Cache[T, PT] {
	return &Cache[T, PT]{db: db}
}

var ErrNilValue = errors.New("nil value pointer")

func (c *Cache[T, PT]) Get(ctx context.Context, key string) (PT, error) {
	r, err := c.db.GetCache(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var zero PT
			return zero, fmt.Errorf("unable to find: %s, err: %w", key, ErrNotFoundInCache)
		}
		var zero PT
		return zero, err
	}
	if time.Until(r.Expires) <= 0 {
		var zero PT
		return zero, fmt.Errorf("key: %s, err: %w", key, ErrOldCache)
	}

	var out T
	ptr := PT(&out)
	_, err = ptr.UnmarshalMsg(r.Data)
	if err != nil {
		var zero PT
		return zero, err
	}

	return ptr, nil
}

func (c *Cache[T, PT]) Set(ctx context.Context, key string, obj PT, exp time.Duration) error {
	if obj == nil {
		return fmt.Errorf("key: %s err: %w", key, ErrNilValue)
	}
	o, err := obj.MarshalMsg(nil)
	if err != nil {
		return err
	}
	return c.db.InsertCache(ctx, db.InsertCacheParams{Key: key, Data: o, Expires: time.Now().Add(exp)})
}
