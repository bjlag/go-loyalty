package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/usecase/login"
)

type Handler struct {
	usecase *login.Usecase
	log     logger.Logger
}

func NewHandler(usecase *login.Usecase, log logger.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if errors.Is(err, errInvalidLogin) || errors.Is(err, errInvalidPassword) {
			h.log.WithError(err).Warning("invalid request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		h.log.WithError(err).Error("error decoding request")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := h.usecase.LoginUser(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, login.ErrUserNotFound) || errors.Is(err, login.ErrWrongPassword) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h.log.WithError(err).Error("error login user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := Response{
		Token: token,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		h.log.WithError(err).Error("could not marshal response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	_, err = w.Write(data)
	if err != nil {
		h.log.WithError(err).Error("could not write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
