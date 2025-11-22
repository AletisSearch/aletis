//go:generate msgp -tests=false

package icons

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/AletisSearch/aletis/internal/cache"
	"github.com/AletisSearch/aletis/internal/db"
	"github.com/go-playground/validator/v10"
	"resty.dev/v3"
)

type Icon struct {
	ContentType string
	IconBytes   []byte
	Expiration  time.Duration
}

type Client struct {
	restyClient *resty.Client
	cache       *cache.Cache[Icon, *Icon]
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func New(q *db.Queries) *Client {
	return &Client{
		restyClient: resty.New(),
		cache:       cache.New[Icon](q),
	}
}

var ErrBadContentType = errors.New("bad content type")

func (c *Client) Get(ctx context.Context, domain string) (*Icon, error) {
	err := validate.Var(domain, "required,fqdn")
	if err != nil {
		return nil, err
	}
	cacheKey := "icon-" + domain
	i, err := c.cache.Get(ctx, cacheKey)
	if err == nil {
		slog.Info("Cache Hit", "Key", cacheKey)
		return i, nil
	}
	if !errors.Is(err, cache.ErrNotFoundInCache) && !errors.Is(err, cache.ErrOldCache) {
		return nil, err
	}
	r, err := c.restyClient.R().
		Get("https://f1.allesedv.com/16/" + domain)
	if err != nil {
		return nil, err
	}
	xs := r.Header().Get("x-source")
	ct := r.Header().Get("content-type")
	if ct == "" {
		return nil, fmt.Errorf("%w: %s", ErrBadContentType, ct)
	}

	i = &Icon{ContentType: ct, IconBytes: r.Bytes()}
	exp := time.Hour * 24 * 7
	if xs != "" {
		exp = time.Minute * 5
	}
	i.Expiration = exp
	return i, c.cache.Set(ctx, cacheKey, i, exp)
}

func (c *Client) Close() error {
	return c.restyClient.Close()
}
