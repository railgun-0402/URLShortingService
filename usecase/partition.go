package usecase

import (
	"context"
	"log"
	"time"
	"url-shorting-service/utils"
)

func tryPartition() {
	ctx := context.Background()

	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// スキーマと親テーブル作成
	if err := utils.EnsureParent(ctx, db); err != nil {
		log.Fatal(err)
	}

	// 例：2025-12〜2026-02 の3ヶ月分パーティション作成
	// （期間集計の pruning を体感するために、複数月にデータをばら撒く）
	if err := utils.EnsureMonthlyPartitions(ctx, db, time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC), 3); err != nil {
		log.Fatal(err)
	}

	// 例：50万件投入（まずこれくらいが体感しやすいらしい）
	n := 500_000
	urlCount := 10_000

	from := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	if err := utils.SeedClickEvents(ctx, db, n, urlCount, from, to); err != nil {
		log.Fatal(err)
	}
}
