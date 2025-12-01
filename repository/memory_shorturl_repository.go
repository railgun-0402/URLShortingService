package repository

import (
	"errors"
	"sync"
	v1 "url-shorting-service/domain/v1"
)

type ShortURLRepository interface {
	Save(s v1.ShortURL) error
	Find(id string) (v1.ShortURL, error)
}

var (
	ErrNotFound      = errors.New("short url not found")
	ErrAlreadyExists = errors.New("id already exists")
)

// InMemoryShortURLRepository とりあえずインメモリ実装（後でDBに差し替える）
type InMemoryShortURLRepository struct {
	mu   sync.RWMutex
	data map[string]v1.ShortURL
}

func NewInMemoryShortURLRepository() *InMemoryShortURLRepository {
	return &InMemoryShortURLRepository{
		data: make(map[string]v1.ShortURL),
	}
}

func (s *InMemoryShortURLRepository) Save(u v1.ShortURL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[u.ID]; ok {
		return ErrAlreadyExists
	}
	s.data[u.ID] = u
	return nil
}

func (s *InMemoryShortURLRepository) Find(id string) (v1.ShortURL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.data[id]
	if !ok {
		return v1.ShortURL{}, ErrNotFound
	}
	return u, nil
}
