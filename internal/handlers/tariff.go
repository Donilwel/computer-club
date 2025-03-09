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

// GetTariff получение полного списка тарифов
func (h tariffHandler) GetTariff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение списка тарифов")
	tariffs, err := h.tariffService.GetTariff(ctx)
	if err != nil {
		h.log.Error("ошибка при запросе списка тарифов")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.log.Info(w, http.StatusOK, "Получен список тарифов")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tariffs)
}

// GetTariffByID получения тарифа по его ID
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
		h.log.Error("ошибка при запросе тарифа по id")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.log.Info(w, http.StatusOK, "Получен тариф по id")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tariff)
}

func NewTariffHandler(tariffService usecase.TariffService, log *logrus.Logger) TariffHandler {
	return &tariffHandler{tariffService: tariffService, log: log}
}
