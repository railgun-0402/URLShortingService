package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"url-shorting-service/handler"
	"url-shorting-service/repository"
	"url-shorting-service/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	repo := repository.NewInMemoryShortURLRepository()
	uc := usecase.NewShortURLUsecase(repo, baseURL)
	h := handler.NewShortURLHandler(uc)
	h.RegisterRoutes(e)

	log.Println("listening on :8080")
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
