package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/repository/models"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
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

func NewComputerHandler(computerService usecase.ComputerService, log *logrus.Logger) ComputerHandler {
	return &computerHandler{
		computerService: computerService,
		log:             log,
	}
}

// GetComputersStatus godoc
// @Summary      Получение статуса компьютеров
// @Description  Возвращает список всех компьютеров и их статусы (доступно только для администраторов)
// @Tags         computers
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Computer "Успешный ответ"
// @Failure      403 {object} string "Ошибка доступа - недостаточно прав"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /computers/status [get]
func (h computerHandler) GetComputersStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение статуса компьютеров")

	role, ok := ctx.Value("role").(string)
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(computers); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}
