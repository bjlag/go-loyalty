package withdraw

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/usecase/withdraw/create"
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

	var req Request

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.WithError(err).Warn("Invalid request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.usecase.CreateWithdraw(r.Context(), userGUID, req.Order, req.Sum)
	if err != nil {
		if errors.Is(err, create.ErrInsufficientBalanceOnAccount) {
			http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}

		h.log.WithError(err).Error("Withdraw error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
