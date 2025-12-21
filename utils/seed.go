package utils

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// とりあえず1回のINSERTに詰める行数を設定
const batchSize = 500

// SeedClickEvents n件のクリックログを投入する
// urlCount: short_url_id の候補数（例：1000なら1..1000をランダムに使う）
func SeedClickEvents(ctx context.Context, db *sql.DB, n int, urlCount int, from time.Time, to time.Time) error {
	// validation
	if err := validationCheck(n, from, to); err != nil {
		return err
	}
	if urlCount <= 0 {
		urlCount = 1
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	start := time.Now()
	inserted := 0

	for inserted < n {
		remain := n - inserted
		bs := batchSize
		if remain < bs {
			bs = remain
		}

		var b strings.Builder
		b.Grow(bs * 64)
		// idはIDENTITY(自動採番)に任せるのでカラムから外す（occurred_at, short_url_id, referrer, user_agent）
		b.WriteString("INSERT INTO click_events (occurred_at, short_url_id, referrer, user_agent) VALUES ")

		args := make([]any, 0, bs*4)
		for i := 0; i < bs; i++ {
			if i > 0 {
				b.WriteString(",")
			}
			// $1, $2...
			base := i*4 + 1
			b.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)", base, base+1, base+2, base+3))

			occurredAt := randomTime(r, from, to)
			shortURLID := int64(r.Intn(urlCount) + 1)

			args = append(args, occurredAt, shortURLID, "https://example.com", "seed-bot/1.0")
		}
		b.WriteString(";")

		if _, err := db.ExecContext(ctx, b.String(), args...); err != nil {
			return fmt.Errorf("insert batch failed (inserted=%d): %w", inserted, err)
		}
		inserted += bs
	}
	fmt.Printf("seeded %d rows in %s\n", inserted, time.Since(start))
	return nil
}

func randomTime(r *rand.Rand, from, to time.Time) time.Time {
	span := to.Unix() - from.Unix()
	if span <= 0 {
		return from
	}
	sec := r.Int63n(span)
	return from.Add(time.Duration(sec) * time.Second)
}

func validationCheck(n int, from time.Time, to time.Time) error {
	fmt.Println("validation_check")

	if n <= 0 {
		return fmt.Errorf("count must be positive: %d", n)
	}
	if !from.Before(to) {
		return fmt.Errorf("from must be before to")
	}
	return nil
}
