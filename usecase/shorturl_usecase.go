package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/url"
	"time"
	"url-shorting-service/domain"
)

const (
	maxIDGenRetry = 5
	shortURLTTL   = 24 * 30 * time.Hour                                              // Expire：30日
	idAlphabet    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 適当なID
)

type ShortURLUsecase struct {
	repo    domain.ShortURLRepository
	baseURL string
}

var (
	id string
	s  domain.ShortURL
)

func NewShortURLUsecase(repo domain.ShortURLRepository, baseURL string) *ShortURLUsecase {
	return &ShortURLUsecase{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (u *ShortURLUsecase) Shorten(ctx context.Context, rawURL string) (domain.ShortURL, error) {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return domain.ShortURL{}, err
	}
	now := time.Now()
	exp := now.Add(shortURLTTL)

	var lastErr error
	for i := 0; i < maxIDGenRetry; i++ {
		tmpID, err := generateID(8)
		if err != nil {
			return domain.ShortURL{}, err
		}
		id = tmpID

		s = domain.ShortURL{
			ID:          id,
			OriginalURL: rawURL,
			CreatedAt:   now,
			ExpiresAt:   &exp,
		}
		if err := u.repo.Save(ctx, s); err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				lastErr = err
				continue // 衝突 → 再トライ
			}
			return domain.ShortURL{}, err
		}
		return s, nil
	}

	if lastErr != nil {
		return domain.ShortURL{}, lastErr
	}
	return domain.ShortURL{}, fmt.Errorf("could not generate unique id")
}

func (u *ShortURLUsecase) Resolve(ctx context.Context, id string) (domain.ShortURL, error) {
	s, err := u.repo.Find(ctx, id)
	if err != nil {
		return domain.ShortURL{}, err
	}
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now()) {
		return domain.ShortURL{}, domain.ErrExpired
	}
	return s, nil
}

// ID生成はとりあえず usecase 側に置く
func generateID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := 0; i < n; i++ {
		b[i] = idAlphabet[int(b[i])%len(idAlphabet)]
	}
	return string(b), nil
}
