package handler

import (
	"errors"
	"net/http"
	"os"
	"url-shorting-service/domain"
	"url-shorting-service/usecase"

	"github.com/labstack/echo/v4"
)

type ShortURLHandler struct {
	uc usecase.ShortURLUseCase
}

func NewShortURLHandler(uc usecase.ShortURLUseCase) *ShortURLHandler {
	return &ShortURLHandler{uc: uc}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"short_url"`
}

func (h *ShortURLHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/shorten", h.Shorten)
	e.GET("/:id", h.Redirect)
}

func (h *ShortURLHandler) Shorten(c echo.Context) error {
	var req shortenRequest
	if err := c.Bind(&req); err != nil || req.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid url")
	}

	s, err := h.uc.Shorten(c.Request().Context(), req.URL)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	resp := shortenResponse{
		ID:       s.ID,
		ShortURL: baseURL + "/" + s.ID,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ShortURLHandler) Redirect(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	s, err := h.uc.Resolve(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.Redirect(http.StatusFound, s.OriginalURL) // 302
}
