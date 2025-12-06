package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"url-shorting-service/domain"
	"url-shorting-service/usecase"

	"github.com/stretchr/testify/require"
)

type mockShortURLRepo struct {
	saveFunc func(ctx context.Context, s domain.ShortURL) error
	findFunc func(ctx context.Context, id string) (domain.ShortURL, error)
}

func (m *mockShortURLRepo) Save(ctx context.Context, s domain.ShortURL) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, s)
	}
	return nil
}

func (m *mockShortURLRepo) Find(ctx context.Context, id string) (domain.ShortURL, error) {
	if m.findFunc != nil {
		return m.findFunc(ctx, id)
	}
	return domain.ShortURL{}, domain.ErrNotFound
}

func TestShorten_Success(t *testing.T) {
	t.Parallel()

	mockRepo := &mockShortURLRepo{
		saveFunc: func(ctx context.Context, s domain.ShortURL) error {
			// 保存される内容を軽くチェック
			if s.ID == "" {
				return errors.New("id is empty")
			}
			if s.OriginalURL != "https://example.com" {
				return errors.New("unexpected original url")
			}
			if s.CreatedAt.IsZero() {
				return errors.New("created_at is zero")
			}
			return nil
		},
	}

	uc := usecase.NewShortURLUsecase(mockRepo, "http://localhost:8080")

	ctx := context.Background()
	s, err := uc.Shorten(ctx, "https://example.com")
	require.NoError(t, err)
	require.Equal(t, "https://example.com", s.OriginalURL)
	require.Len(t, s.ID, 8)
}

func TestShorten_InvalidURL(t *testing.T) {
	t.Parallel()

	mockRepo := &mockShortURLRepo{}
	uc := usecase.NewShortURLUsecase(mockRepo, "http://localhost:8080")

	_, err := uc.Shorten(context.Background(), "hogehoge")
	require.Error(t, err, "invalid URL はエラーになるべき")
}

// 衝突時にリトライされるか
func TestShorten_RetryOnCollision(t *testing.T) {
	t.Parallel()

	callCount := 0

	mockRepo := &mockShortURLRepo{
		saveFunc: func(ctx context.Context, s domain.ShortURL) error {
			callCount++
			if callCount == 1 {
				// 1回目は衝突扱い
				return domain.ErrAlreadyExists
			}
			// 2回目以降は成功
			return nil
		},
	}

	uc := usecase.NewShortURLUsecase(mockRepo, "http://localhost:8080")

	_, err := uc.Shorten(context.Background(), "https://example.com")
	require.NoError(t, err)
	require.Equal(t, 2, callCount, "衝突後にリトライされているはず")
}

// ---- Resolve のテスト ----

func TestResolve_Success(t *testing.T) {
	t.Parallel()

	expected := domain.ShortURL{
		ID:          "abc12345",
		OriginalURL: "https://example.com",
		CreatedAt:   time.Now(),
	}

	mockRepo := &mockShortURLRepo{
		findFunc: func(ctx context.Context, id string) (domain.ShortURL, error) {
			if id != "abc12345" {
				return domain.ShortURL{}, domain.ErrNotFound
			}
			return expected, nil
		},
	}

	uc := usecase.NewShortURLUsecase(mockRepo, "http://localhost:8080")

	got, err := uc.Resolve(context.Background(), "abc12345")
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestResolve_NotFound(t *testing.T) {
	t.Parallel()

	mockRepo := &mockShortURLRepo{
		findFunc: func(ctx context.Context, id string) (domain.ShortURL, error) {
			return domain.ShortURL{}, domain.ErrNotFound
		},
	}

	uc := usecase.NewShortURLUsecase(mockRepo, "http://localhost:8080")

	_, err := uc.Resolve(context.Background(), "notfound")
	require.ErrorIs(t, err, domain.ErrNotFound)
}
