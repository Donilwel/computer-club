package httpService

import (
	"computer-club/internal/usecase"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	userService     usecase.UserService
	computerService usecase.ComputerService
	sessionService  usecase.SessionService
	tariffService   usecase.TariffService
	log             *logrus.Logger
}

// NewHandler создает новый HTTP-обработчик
func NewHandler(userService usecase.UserService,
	computerService usecase.ComputerService,
	sessionService usecase.SessionService,
	tariffService usecase.TariffService,
	log *logrus.Logger) *Handler {
	return &Handler{
		userService:     userService,
		computerService: computerService,
		sessionService:  sessionService,
		tariffService:   tariffService,
		log:             log,
	}
}
