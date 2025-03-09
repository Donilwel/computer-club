package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// GetTariffByID получения тарифа по его ID
func (h *Handler) GetTariffByID(w http.ResponseWriter, r *http.Request) {
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
