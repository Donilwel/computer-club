package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"computer-club/internal/models"
	"encoding/json"
	"net/http"
)

func (h *Handler) PutMoneyOnWallet(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на отправку средств на счет игрока")

	role, ok := r.Context().Value("role").(string)
	if !ok || role != string(models.Admin) {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при переводе: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()

	if err := h.walletService.Deposit(req.UserID, req.Amount); err != nil {
		h.log.WithError(err).Error("Ошибка при передаче денег")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	transaction, err := h.walletService.CreateTransaction(req.UserID, req.Amount, string(models.Add), -1)
	if err != nil {
		h.log.WithError(err).Error("Ошибка при создании модели транзакции")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
