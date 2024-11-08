package upload

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/model"
	"github.com/bjlag/go-loyalty/internal/usecase/accrual/create"
)

type Handler struct {
	usecase *create.Usecase
	log     logger.Logger
}

func NewHandler(usecase *create.Usecase, log logger.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
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

	orderNumber := b.String()

	if err := h.usecase.CreateAccrual(ctx, model.NewAccrual(orderNumber, userGUID)); err != nil {
		switch {
		case errors.Is(err, create.ErrInvalidOrderNumber):
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		case errors.Is(err, create.ErrAnotherUserHasAlreadyRegisteredOrder):
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		case errors.Is(err, create.ErrOrderAlreadyExists):
			w.WriteHeader(http.StatusOK)
			return
		}

		h.log.WithError(err).Error("Error creating accrual")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
