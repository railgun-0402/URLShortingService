package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/url"
	"time"
	domain "url-shorting-service/domain/v1"
	"url-shorting-service/repository"
)

// 適当なID
const idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ShortURLUseCase interface {
	Shorten(ctx context.Context, rawURL string) (domain.ShortURL, error)
	Resolve(ctx context.Context, id string) (domain.ShortURL, error)
}

type shortURLUsecase struct {
	repo    repository.ShortURLRepository
	baseURL string
}

func NewShortURLUsecase(repo repository.ShortURLRepository, baseURL string) ShortURLUseCase {
	return &shortURLUsecase{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (u *shortURLUsecase) Shorten(ctx context.Context, rawURL string) (domain.ShortURL, error) {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return domain.ShortURL{}, err
	}

	var (
		id string
		s  domain.ShortURL
	)

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
		if err := u.repo.Save(s); err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				continue // 衝突 → 再トライ
			}
			return domain.ShortURL{}, err
		}
		return s, nil
	}

	return domain.ShortURL{}, fmt.Errorf("could not generate unique id")
}

func (u *shortURLUsecase) Resolve(ctx context.Context, id string) (domain.ShortURL, error) {
	return u.repo.Find(id)
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
