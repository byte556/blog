package models

import "errors"

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameIsEmpty = errors.New("username is empty")
	ErrPasswordIsEmpty = errors.New("password is empty")
	ErrConvert         = errors.New("convert error")
)
