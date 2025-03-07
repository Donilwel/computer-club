package httpService

import (
	"computer-club/internal/middleware"
	"encoding/json"
	"net/http"
)

func (h *Handler) GetTariff(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на получение списка тарифов")
	tariffs, err := h.tariffService.GetTariff()
	if err != nil {
		h.log.Error("ошибка при запросе списка тарифов")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.log.Info(w, http.StatusOK, "Получен список тарифов")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tariffs)
}
