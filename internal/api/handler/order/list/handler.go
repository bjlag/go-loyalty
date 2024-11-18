package list

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
)

type Handler struct {
	repo repository.AccrualRepo
	log  logger.Logger
}

func NewHandler(repo repository.AccrualRepo, log logger.Logger) *Handler {
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

	rows, err := h.repo.AccrualsByUser(ctx, userGUID)
	if err != nil {
		h.log.WithError(err).Error("Could not get accruals")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if rows == nil {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

	resp := make(Response, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, Order{
			Number:     row.OrderNumber,
			Status:     strings.ToUpper(row.Status.String()),
			Accrual:    row.Accrual,
			UploadedAt: row.UploadedAt,
		})
	}

	// todo дата в формате 2020-12-10T15:15:45+03:00
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
