package httpService

import (
	"computer-club/internal/usecase"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	userService     usecase.UserService
	computerService usecase.ComputerService
	sessionService  usecase.SessionService
	log             *logrus.Logger
}

// NewHandler создает новый HTTP-обработчик
func NewHandler(userService usecase.UserService,
	computerService usecase.ComputerService,
	sessionService usecase.SessionService,
	log *logrus.Logger) *Handler {
	return &Handler{
		userService:     userService,
		computerService: computerService,
		sessionService:  sessionService,
		log:             log,
	}
}
