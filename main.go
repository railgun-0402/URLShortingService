package main

import (
	"crypto/rand"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type ShortURL struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
}

var (
	ErrNotFound      = errors.New("short url not found")
	ErrAlreadyExists = errors.New("id already exists")
)

type URLStore interface {
	Save(s ShortURL) error
	Find(id string) (ShortURL, error)
}

// InMemoryURLStore とりあえずインメモリ実装（後でDBに差し替える）
type InMemoryURLStore struct {
	mu   sync.RWMutex
	data map[string]ShortURL
}

func NewInMemoryURLStore() *InMemoryURLStore {
	return &InMemoryURLStore{
		data: make(map[string]ShortURL),
	}
}

func (s *InMemoryURLStore) Save(u ShortURL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[u.ID]; ok {
		return ErrAlreadyExists
	}
	s.data[u.ID] = u
	return nil
}

func (s *InMemoryURLStore) Find(id string) (ShortURL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.data[id]
	if !ok {
		return ShortURL{}, ErrNotFound
	}
	return u, nil
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"short_url"`
}

// 適当なID
const idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

func main() {
	e := echo.New()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	store := NewInMemoryURLStore()

	e.POST("/shorten", func(c echo.Context) error {
		var req ShortenRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if req.URL == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "url is required")
		}

		// validation check
		if _, err := url.ParseRequestURI(req.URL); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid url format")
		}

		var id string
		// ID生成 & 衝突確認
		// 一旦はインメモリに設定した内容を見る(後にDBも加える)
		for i := 0; i < 5; i++ {
			tmpID, err := generateID(8)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate id")
			}
			id = tmpID
			err = store.Save(ShortURL{
				ID:          id,
				OriginalURL: req.URL,
				CreatedAt:   time.Now(),
			})
			if err == nil {
				break
			}
			if !errors.Is(err, ErrAlreadyExists) {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to save url")
			}
		}
		if id == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not generate unique id")
		}

		resp := ShortenResponse{
			ID:       id,
			ShortURL: baseURL + "/" + id,
		}
		return c.JSON(http.StatusOK, resp)
	})

	// リダイレクト用エンドポイント
	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "id is required")
		}

		u, err := store.Find(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "short url not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to lookup url")
		}

		// 302でリダイレクト（v1ではこれを採用）
		return c.Redirect(http.StatusFound, u.OriginalURL)
	})

	log.Println("listening on :8080")
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
