package models

import "errors"

var (
	ErrUserNotExist  = errors.New("user not exist")
	ErrInvalidPasswd = errors.New("Passwd or username not right")
	ErrInvalidParams = errors.New("Invalid params")
	ErrNicknameExist = errors.New("nickname exist")
)
