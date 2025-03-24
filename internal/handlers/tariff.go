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

type TariffHandler interface {
	GetTariff(w http.ResponseWriter, r *http.Request)
	GetTariffByID(w http.ResponseWriter, r *http.Request)
	ChangeTariff(w http.ResponseWriter, r *http.Request)
	AddTariff(w http.ResponseWriter, r *http.Request)
	DeleteTariff(w http.ResponseWriter, r *http.Request)
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
func (h *tariffHandler) GetTariff(w http.ResponseWriter, r *http.Request) {
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
		return
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
func (h *tariffHandler) GetTariffByID(w http.ResponseWriter, r *http.Request) {
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

type TariffBody struct {
	TariffID int64   `json:"tariff_id" example:"123"`
	Name     string  `json:"name" example:"Free"`
	Price    float64 `json:"price" example:"0"`
	Duration int64   `json:"duration" example:"10"`
}

func (h *tariffHandler) ChangeTariff(w http.ResponseWriter, r *http.Request) {

}

// AddTariff godoc
// @Summary      Создание нового тарифа
// @Description  Создает новый тариф с указанными данными (доступен для администратора)
// @Tags         tariffs
// @Accept       json
// @Produce      json
// @Param        request body TariffBody true "Данные для тарифа"
// @Success      200 {object} models.Tariff "Тариф успешно создан"
// @Failure      403 {object} string "Недостаточно прав: нужен пользователь с ролью админ"
// @Failure      400 {object} string "Некорректные данные (пустое имя, отрицательный айди или цена или длительность)"
// @Failure      409 {object} string "Айди тарифа уже используется"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /tariff [post]
func (h *tariffHandler) AddTariff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Добавление тарифа в базу данных")
	var req TariffBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.WithError(err).Error("ошибка чтения тела JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidTariffID.Error())
		return
	}

	defer r.Body.Close()
	role, ok := ctx.Value("role").(string)
	if !ok || role != string(models.Admin) {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при добавления тарифа: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	tariff, err := h.tariffService.CreateTariff(ctx, req.TariffID, req.Name, req.Price, req.Duration)
	if err != nil {
		switch err {
		case errors.ErrInvalidInputDataTariff:
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.ErrTariffWithIdAlreadyExist:
			middleware.WriteError(w, http.StatusConflict, err.Error())
		default:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		h.log.Error(err)
		return
	}
	h.log.Info("Тариф успешно создан")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tariff); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
		return
	}
}

// DeleteTariff godoc
// @Summary удаление Тарифа
// @Description удаление Тарифа из базы данных (доступно для админа)
// @Tags tariffs
// @Produce      json
// @Security     BearerAuth
// @Success      200  string message": "тариф успешно удален"
// @Failure		 400 {object} string "Некорректный ID тарифа"
// @Failure      403 {object} string "Ошибка доступа - недостаточно прав"
// @Failure		 404 {object} string "Тариф не найден"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /tariff/{id} [delete]
func (h *tariffHandler) DeleteTariff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Удаление тарифа из базы данных")
	role, ok := ctx.Value("role").(string)
	if !ok || role != string(models.Admin) {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при удалении тарифа: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.log.Error("ID тарифа не должен быть пустым")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidTariffID.Error())
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.WithError(err).Error("Некорректный ID тарифа")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidTariffID.Error())
		return
	}

	if id <= 0 {
		h.log.WithError(err).Error("Некорректный ID тарифа")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidTariffID.Error())
		return
	}

	if err := h.tariffService.DeleteTariffByID(ctx, id); err != nil {
		switch err {
		case errors.ErrFindTariffByID:
			middleware.WriteError(w, http.StatusNotFound, err.Error())
		default:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		h.log.WithError(err).Error("ошибка удаления тарифа")
		return
	}

	h.log.Info("тариф успешно удален")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "тариф успешно удален"}); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
		return
	}
}
