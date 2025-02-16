package domain

import "errors"

type HttpError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

var (
	ErrBadRequest              = errors.New("some troubles in request")
	ErrInternalServerError     = errors.New("something went wrong :(")
	ErrUserAlreadyExists       = errors.New("client already exists")
	ErrAdvertiserAlreadyExists = errors.New("advertiser already exists")
	ErrNotFound                = errors.New("not found")
	ErrUserNotFound            = errors.New("client not found")
	ErrAdvertiserNotFound      = errors.New("advertiser not found")

	ErrNewDateLowerThanCurrent = errors.New("new date must be bigger than current")
	ErrModerationNotPassed     = errors.New("moderation not passed")
)
