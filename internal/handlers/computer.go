package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/models"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ComputerHandler interface {
	GetComputersStatus(w http.ResponseWriter, r *http.Request)
	AddComputer(w http.ResponseWriter, r *http.Request)
	DeleteComputer(w http.ResponseWriter, r *http.Request)
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
// @Router       /computers [get]
func (h *computerHandler) GetComputersStatus(w http.ResponseWriter, r *http.Request) {
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

// DeleteComputer godoc
// @Summary      Удаляет компьютер
// @Description  Удаляет компьютер из базы данных (доступно только для администраторов)
// @Tags         computers
// @Produce      json
// @Security     BearerAuth
// @Success      200  string message": "Компьютер успешно удален
// @Failure		 400 {object} string "Некорректный ID компьютера"
// @Failure      403 {object} string "Ошибка доступа - недостаточно прав"
// @Failure		 404 {object} string "Компьютер не найден"
// @Failure		 409 {object} string "Компьютер занят"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /computers/{id} [delete]
func (h *computerHandler) DeleteComputer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role, ok := ctx.Value("role").(string)
	if !ok || role != string(models.Admin) {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при удалении компьютера: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.WithError(err).Error("Некорректный ID компьютера")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidTariffID.Error())
		return
	}

	err = h.computerService.DeleteComputer(ctx, int(id))
	if err != nil {
		switch err {
		case errors.ErrFindComputer:
			middleware.WriteError(w, http.StatusNotFound, errors.ErrFindComputer.Error())
		case errors.ErrDeleteComputer:
			middleware.WriteError(w, http.StatusInternalServerError, errors.ErrDeleteComputer.Error())
		case errors.ErrPCBusy:
			middleware.WriteError(w, http.StatusConflict, errors.ErrPCBusy.Error())
		default:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		h.log.WithError(err).Error("Ошибка удаления компьютера")
		return
	}
	h.log.Info("Компьютер удален успешно")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Компьютер успешно удален"}); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}

// AddComputer godoc
// @Summary      Добавляет компьютер
// @Description  Добавляет компьютер в базу данных (доступно только для администраторов)
// @Tags         computers
// @Produce      json
// @Security     BearerAuth
// @Success      200  string message": "Компьютер успешно добавлен
// @Failure      403 {object} string "Ошибка доступа - недостаточно прав"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /computers [post]
func (h *computerHandler) AddComputer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Добавление компьютера")
	role, ok := ctx.Value("role").(string)
	if !ok || role != string(models.Admin) {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при добавлении компьютера: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	computer, err := h.computerService.AddComputer(ctx)
	if err != nil {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при добавлении компьютера")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCreateComputer.Error())
		return
	}
	h.log.Info("Компьютер добавлен успешно")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(computer); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}
