package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type TariffHandler interface {
	GetTariff(w http.ResponseWriter, r *http.Request)
	GetTariffByID(w http.ResponseWriter, r *http.Request)
}

type tariffHandler struct {
	tariffService usecase.TariffService
	log           *logrus.Logger
}

func NewTariffHandler(tariffService usecase.TariffService, log *logrus.Logger) TariffHandler {
	return &tariffHandler{tariffService: tariffService, log: log}
}

// GetTariff godoc
// @Summary      Получение списка тарифов
// @Description  Возвращает список всех доступных тарифов
// @Tags         tariffs
// @Produce      json
// @Success      200 {array} models.Tariff
// @Failure      500 {object} string "ошибка при поиске тарифов в базе данных"
// @Router       /tariff [get]
func (h tariffHandler) GetTariff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение списка тарифов")
	tariffs, err := h.tariffService.GetTariff(ctx)
	if err != nil {
		h.log.Error(err)
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.log.Info(w, http.StatusOK, "Получен список тарифов")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tariffs); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}

// GetTariffByID godoc
// @Summary      Получение тарифа по ID
// @Description  Возвращает информацию о тарифе по указанному ID
// @Tags         tariffs
// @Produce      json
// @Param        id path int true "ID тарифа"
// @Success      200 {object} models.Tariff
// @Failure      400 {object} string "ошибка чтения тариф id"
// @Failure      404 {object} string "тариф не найден"
// @Failure      500 {object} string "внутренняя ошибка сервера"
// @Router       /tariff/{id} [get]
func (h tariffHandler) GetTariffByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.WithError(err).Error("Некорректный ID тарифа")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidTariffID.Error())
		return
	}

	h.log.Info("Запрос на получение тарифа по id")
	tariff, err := h.tariffService.GetTariffByID(ctx, id)
	if err != nil {
		h.log.WithError(err).Error("ошибка при запросе тарифа по id")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.log.Info(w, http.StatusOK, "Получен тариф по id")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tariff); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}
