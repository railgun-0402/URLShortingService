package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"url-shorting-service/handler"
	"url-shorting-service/repository/postgres"
	"url-shorting-service/usecase"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// local
		dsn = "postgres://urlshort:urlshort@localhost:5432/urlshort?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// v1
	// repo := repository.NewInMemoryShortURLRepository()

	// v2
	repo := postgres.NewPostgresShortURLRepository(db)
	uc := usecase.NewShortURLUsecase(repo, baseURL)
	h := handler.NewShortURLHandler(uc)
	h.RegisterRoutes(e)

	log.Println("listening on :8080")
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
