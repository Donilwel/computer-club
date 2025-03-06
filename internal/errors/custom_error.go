package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidRole        = errors.New("некорректная роль")
	ErrUserAlreadyExists  = errors.New("пользователь уже существует")
	ErrInvalidCredentials = errors.New("неверный email или пароль")
	ErrUsernameTaken      = errors.New("пользователь с таким никнеймом уже существует")

	ErrSessionNotFound  = errors.New("сессия не найдена")
	ErrSessionActive    = errors.New("у пользователя уже есть активная сессия")
	ErrPCBusy           = errors.New("компьютер уже занят")
	ErrInvalidSessionID = errors.New("некорректный идентификатор сессии")
	ErrRegistration     = errors.New("ошибка при регистрации пользователя")
	ErrComputerNotFound = errors.New("компьютер не найден")
	ErrTokenGeneration  = errors.New("ошибка генерации токена")
)
