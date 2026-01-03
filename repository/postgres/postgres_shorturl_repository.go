package postgres

import (
	"context"
	"database/sql"
	"errors"
	"url-shorting-service/domain"
)

type ShortURLRepository struct {
	db *sql.DB
}

func NewPostgresShortURLRepository(db *sql.DB) *ShortURLRepository {
	return &ShortURLRepository{db: db}
}

func (r *ShortURLRepository) Save(ctx context.Context, s domain.ShortURL) error {
	// ON CONFLICT で衝突時は何もしない → rowsAffected = 0 なら ErrAlreadyExists とみなす
	const q = `
		INSERT INTO short_urls (id, original_url, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING;
		`
	res, err := r.db.ExecContext(ctx, q, s.ID, s.OriginalURL, s.CreatedAt, s.ExpiresAt)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrAlreadyExists
	}
	return nil
}

func (r *ShortURLRepository) Find(ctx context.Context, id string) (domain.ShortURL, error) {
	const q = `SELECT id, original_url, created_at, expires_at FROM short_urls WHERE id = $1;`

	var s domain.ShortURL
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&s.ID,
		&s.OriginalURL,
		&s.CreatedAt,
		&s.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ShortURL{}, domain.ErrNotFound
		}
		return domain.ShortURL{}, err
	}
	return s, nil
}
