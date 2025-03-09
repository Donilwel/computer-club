package handlers

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"computer-club/internal/models"
	"computer-club/internal/usecase"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ComputerHandler interface {
	GetComputersStatus(w http.ResponseWriter, r *http.Request)
}

type computerHandler struct {
	computerService usecase.ComputerService
	log             *logrus.Logger
}

// GetComputersStatus возвращает статус компьютеров
func (h computerHandler) GetComputersStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение статуса компьютеров")

	role, ok := r.Context().Value("role").(string)
	if !ok || role != string(models.Admin) {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при получении списка компьютеров: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	computers, err := h.computerService.GetComputersStatus(ctx)
	if err != nil {
		h.log.WithError(err).Error("Ошибка при получении списка компьютеров")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.WithField("count", len(computers)).Info("Статус компьютеров получен")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(computers)
}

func NewComputerHandler(computerService usecase.ComputerService, log *logrus.Logger) ComputerHandler {
	return &computerHandler{
		computerService: computerService,
		log:             log,
	}
}
