package domain

import (
	"context"
)

type ShortURLRepository interface {
	Save(ctx context.Context, s ShortURL) error
	Find(ctx context.Context, id string) (ShortURL, error)
}
