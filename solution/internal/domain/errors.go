package domain

import "errors"

type HttpError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

var (
	ErrBadRequest        = errors.New("some troubles in request")
	ErrInternalServerError = errors.New("something went wrong :(")
	ErrUserAlreadyExists = errors.New("client already exists")
	ErrAdvertiserAlreadyExists = errors.New("advertiser already exists")
	ErrUserNotFound      = errors.New("client not found")
	ErrAdvertiserNotFound = errors.New("advertiser not found")
)
