package application

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidEmail     = errors.New("Некорректный формат адреса электронной почты")
	ErrEmailAlreadyUsed = errors.New("Адрес электронной почты уже используется")
)

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}
