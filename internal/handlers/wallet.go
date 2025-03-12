package handlers

import (
	"computer-club/internal/middleware"
	models2 "computer-club/internal/repository/models"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type WalletHandler interface {
	PutMoneyOnWallet(http.ResponseWriter, *http.Request)
}

func NewWalletHandler(walletService usecase.WalletService, log *logrus.Logger) WalletHandler {
	return &walletHandler{walletService: walletService, log: log}
}

type walletHandler struct {
	walletService usecase.WalletService
	log           *logrus.Logger
}

func (h walletHandler) PutMoneyOnWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на отправку средств на счет игрока")

	role, ok := r.Context().Value("role").(string)
	if !ok || role != string(models2.Admin) {
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

	transaction, err := h.walletService.PutMoneyOnWallet(ctx, req.UserID, req.Amount)
	if err != nil {
		h.log.WithError(err).Error("Ошибка при переводе денег или создании модели транзакции")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}
