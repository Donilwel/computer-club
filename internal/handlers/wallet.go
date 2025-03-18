package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/models"
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

type PutMoneyReq struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

// PutMoneyOnWallet godoc
// @Summary      Пополнение счета игрока
// @Description  Администратор пополняет счет указанного пользователя (admin only)
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body PutMoneyReq true "Данные для пополнения счета"
// @Success      200 {object} models.Transaction "Счет успешно пополнен"
// @Failure      400 {object} string "Некорректный запрос (например, сумма ≤ 0)"
// @Failure      403 {object} string "Ошибка доступа - недостаточно прав"
// @Failure      404 {object} string "Пользователь не найден"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /pay [put]
func (h walletHandler) PutMoneyOnWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на отправку средств на счет игрока")

	role, ok := r.Context().Value("role").(string)
	if !ok || role != string(models.Admin) {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при переводе: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	var req PutMoneyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()

	transaction, err := h.walletService.PutMoneyOnWallet(ctx, req.UserID, req.Amount)
	if err != nil {
		switch err {
		case errors.ErrFindUser, errors.ErrCheckBalance:
			middleware.WriteError(w, http.StatusNotFound, err.Error())
		case errors.ErrInvalidAmount:
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
		default:

			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		h.log.WithError(err).Error("Ошибка при переводе денег или создании модели транзакции")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}
