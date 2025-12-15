package domain

import "errors"

var (
	ErrNotFound      = errors.New("short url not found")
	ErrAlreadyExists = errors.New("id already exists")
	ErrExpired       = errors.New("short url expired")
)
