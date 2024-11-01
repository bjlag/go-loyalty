package register

import (
	"encoding/json"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/usecase/register"
)

type Handler struct {
	usecase *register.Usecase
	log     logger.Logger
}

func NewHandler(usecase *register.Usecase, log logger.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.WithError(err).Error("error decoding request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := h.usecase.RegisterUser(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.WithError(err).Error("error registering user")
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

	_, err = w.Write(data)
	if err != nil {
		h.log.WithError(err).Error("could not write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}