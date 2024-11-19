package get

import (
	"encoding/json"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
)

type Handler struct {
	repo repository.AccountRepo
	log  logger.Logger
}

func NewHandler(repo repository.AccountRepo, log logger.Logger) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
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

	balance, err := h.repo.Balance(ctx, userGUID)
	if err != nil {
		h.log.WithError(err).Error("Could not get balance")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := &Response{
		Current:   balance,
		Withdrawn: 0,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.WithError(err).Error("Could not write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
