package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"encoding/json"
	"net/http"
)

// GetComputersStatus возвращает статус компьютеров
func (h *Handler) GetComputersStatus(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на получение статуса компьютеров")

	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при получении списка компьютеров: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	computers, err := h.computerService.GetComputersStatus()
	if err != nil {
		h.log.WithError(err).Error("Ошибка при получении списка компьютеров")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.WithField("count", len(computers)).Info("Статус компьютеров получен")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(computers)
}
