package storage

import "errors"

var (
	ErrUserExists             = errors.New("user already exist")
	ErrUserNotFound           = errors.New("user not found")
	ErrAppNotFound            = errors.New("app not found")
	ErrPermissionAlreadyExist = errors.New("permission already exist")
	ErrNoPermission           = errors.New("user don't have this permission")
)
