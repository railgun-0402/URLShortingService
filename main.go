package main

import (
	"context"
	"log"
	"time"
	"url-shorting-service/utils"

	_ "github.com/lib/pq"
)

func main() {
	tryPartition()

	//e := echo.New()
	//
	//baseURL := os.Getenv("BASE_URL")
	//if baseURL == "" {
	//	baseURL = "http://localhost:8080"
	//}
	//
	//dsn := os.Getenv("DATABASE_URL")
	//if dsn == "" {
	//	// local
	//	dsn = "postgres://urlshort:urlshort@localhost:5432/urlshort?sslmode=disable"
	//}
	//
	//db, err := sql.Open("postgres", dsn)
	//if err != nil {
	//	log.Fatalf("failed to open db: %v", err)
	//}
	//defer db.Close()
	//
	//if err := db.Ping(); err != nil {
	//	log.Fatalf("failed to ping db: %v", err)
	//}
	//
	//// v1
	//// repo := repository.NewInMemoryShortURLRepository()
	//
	//// v2
	//repo := postgres.NewPostgresShortURLRepository(db)
	//uc := usecase.NewShortURLUsecase(repo, baseURL)
	//h := handler.NewShortURLHandler(uc)
	//h.RegisterRoutes(e)
	//
	//log.Println("listening on :8080")
	//if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//	log.Fatalf("server error: %v", err)
	//}
}

// 後で整理する
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
