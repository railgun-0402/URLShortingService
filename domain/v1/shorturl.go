package v1

import "time"

type ShortURL struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
}
