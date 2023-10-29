package repo

import (
	"context"
	"errors"
	"fmt"
	"server/internal/infra"
)

const KeyNotFound = "key not found"

type CrawlerRepo interface {
	GetUrl(ctx context.Context, url string) (string, error)
	StoreUrl(ctx context.Context, key, value string) error
}

type crawlerRepository struct {
	client infra.RedisClient
}

// NewCrawlerRepository builds a crawlerRepository and injects its dependencies
func NewCrawlerRepository(client infra.RedisClient) CrawlerRepo {
	return &crawlerRepository{
		client: client,
	}
}

func (r *crawlerRepository) GetUrl(ctx context.Context, url string) (string, error) {
	url, err := r.client.Get(ctx, url)
	if err != nil {
		if err.Error() == infra.RedisKeyNotFound {
			return "", errors.New(KeyNotFound)
		}

		return "", errors.New(fmt.Sprintf("error getting url from repo: %s", err))
	}

	return url, nil

}

func (r *crawlerRepository) StoreUrl(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value)
}
