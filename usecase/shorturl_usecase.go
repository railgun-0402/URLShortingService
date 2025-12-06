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

// 適当なID
const idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

	// 5回は任意の回数で試してるだけ
	for i := 0; i < 5; i++ {
		tmpID, err := generateID(8)
		if err != nil {
			return domain.ShortURL{}, err
		}
		id = tmpID

		s = domain.ShortURL{
			ID:          id,
			OriginalURL: rawURL,
			CreatedAt:   time.Now(),
		}
		if err := u.repo.Save(ctx, s); err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				continue // 衝突 → 再トライ
			}
			return domain.ShortURL{}, err
		}
		return s, nil
	}

	return domain.ShortURL{}, fmt.Errorf("could not generate unique id")
}

func (u *ShortURLUsecase) Resolve(ctx context.Context, id string) (domain.ShortURL, error) {
	return u.repo.Find(ctx, id)
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
