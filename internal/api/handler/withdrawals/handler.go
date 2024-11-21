package withdrawals

import (
	"encoding/json"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/api"
	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
)

type Handler struct {
	repo repository.TransactionRepo
	log  logger.Logger
}

func NewHandler(repo repository.TransactionRepo, log logger.Logger) *Handler {
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

	rows, err := h.repo.Withdrawals(ctx, userGUID)
	if err != nil {
		h.log.WithError(err).Error("Could not get withdrawals")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if rows == nil {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

	resp := make(Response, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, Withdraw{
			Order:       row.OrderNumber,
			Sum:         float32(row.Sum),
			ProcessedAt: api.Datetime(row.ProcessedAt),
		})
	}

	data, err := json.Marshal(resp)
	if err != nil {
		h.log.WithError(err).Error("Could not marshal response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(data)
	if err != nil {
		h.log.WithError(err).Error("Could not write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
