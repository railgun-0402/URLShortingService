package repository

import (
	"sync"
	"url-shorting-service/domain"
)

type ShortURLRepository interface {
	Save(s domain.ShortURL) error
	Find(id string) (domain.ShortURL, error)
}

// InMemoryShortURLRepository とりあえずインメモリ実装（後でDBに差し替える）
type InMemoryShortURLRepository struct {
	mu   sync.RWMutex
	data map[string]domain.ShortURL
}

func NewInMemoryShortURLRepository() *InMemoryShortURLRepository {
	return &InMemoryShortURLRepository{
		data: make(map[string]domain.ShortURL),
	}
}

func (s *InMemoryShortURLRepository) Save(u domain.ShortURL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[u.ID]; ok {
		return domain.ErrAlreadyExists
	}
	s.data[u.ID] = u
	return nil
}

func (s *InMemoryShortURLRepository) Find(id string) (domain.ShortURL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.data[id]
	if !ok {
		return domain.ShortURL{}, domain.ErrNotFound
	}
	return u, nil
}
