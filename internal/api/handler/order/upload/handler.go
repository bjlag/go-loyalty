package upload

import (
	"bytes"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	"github.com/bjlag/go-loyalty/internal/infrastructure/validator"
	"github.com/bjlag/go-loyalty/internal/model"
)

type Handler struct {
	jwt *auth.JWTBuilder
	rep repository.AccrualRepository
	log logger.Logger
}

func NewHandler(jwt *auth.JWTBuilder, rep repository.AccrualRepository, log logger.Logger) *Handler {
	return &Handler{
		jwt: jwt,
		rep: rep,
		log: log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userGUID, err := auth.UserGUIDFromContext(ctx)
	if err != nil {
		h.log.WithError(err).Error("Could not get user GUID from context")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var b bytes.Buffer

	_, err = b.ReadFrom(r.Body)
	if err != nil {
		h.log.WithError(err).Error("Error reading body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	number := b.String()
	if !validator.CheckLuhn(number) {
		h.log.WithField("number", number).Warning("Order number is invalid")
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	// todo проверить есть ли у этого пользователя этот заказ в обработке 200
	// todo проверить если ли этот заказ в обработке в принципе 409

	accrual := model.NewAccrual(number, userGUID)
	err = h.rep.Insert(ctx, accrual)
	if err != nil {
		h.log.WithError(err).Error("Error inserting accrual")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// todo асинхронное получение балов лояльности по заказу

	w.WriteHeader(http.StatusAccepted)
}
