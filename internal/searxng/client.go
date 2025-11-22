package searxng

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/AletisSearch/aletis/internal/cache"
	"github.com/AletisSearch/aletis/internal/db"
	"resty.dev/v3"
)

type Client struct {
	url         string
	restyClient *resty.Client
	cache       *cache.Cache[SearchResponse, *SearchResponse]
}

func NewClient(url string, q *db.Queries) *Client {
	return &Client{
		url: url,
		restyClient: resty.New().
			SetHeader("Accept", "application/json, text/html").
			SetHeader("Accept-Language", "*").SetBaseURL(url),
		cache: cache.New[SearchResponse](q),
	}
}

func (c *Client) Search(ctx context.Context, query string, page ...int) (*SearchResponse, error) {
	pageno := "1"
	if len(page) != 0 {
		pageno = strconv.Itoa(page[0])
	}

	cacheKey := "search-" + query + "-" + pageno
	s, err := c.cache.Get(ctx, cacheKey)
	if err == nil {
		slog.Info("Cache Hit", "Key", cacheKey)
		return s, nil
	}
	if !errors.Is(err, cache.ErrNotFoundInCache) && !errors.Is(err, cache.ErrOldCache) {
		return nil, err
	}
	queryParams := map[string]string{
		"q":      query,
		"format": "json",
	}
	if len(page) != 0 {
		queryParams["pageno"] = pageno
	}

	res, err := c.restyClient.R().WithContext(ctx).
		SetQueryParams(queryParams).
		Get("/search")
	if err != nil || res.StatusCode() >= 400 {
		return nil, fmt.Errorf("unable to fetch searxng response status: %d error: %w", res.StatusCode(), err)
	}
	defer res.Body.Close()

	rawJSON := res.Bytes()
	var sr SearchResponse
	if err = json.Unmarshal(rawJSON, &sr); err != nil {
		return nil, fmt.Errorf("unable to unmarshal searxng response error: %w\nraw JSON:\n%s\n", err, rawJSON)
	}
	OrderResults(&sr.Results)

	return &sr, c.cache.Set(ctx, cacheKey, &sr, time.Minute*15)
}

func (c *Client) Close() {
	c.restyClient.Close()
}
