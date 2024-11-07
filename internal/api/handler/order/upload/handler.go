package upload

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
)

type Handler struct {
	jwt *auth.JWTBuilder
	log logger.Logger
}

func NewHandler(jwt *auth.JWTBuilder, log logger.Logger) *Handler {
	return &Handler{
		jwt: jwt,
		log: log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userGUID, err := auth.UserGUIDFromContext(r.Context())
	if err != nil {
		h.log.WithError(err).Error("could not get user GUID from context")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, userGUID)
	fmt.Fprint(w, "\n")

	var b bytes.Buffer

	_, err = b.ReadFrom(r.Body)
	if err != nil {
		h.log.WithError(err).Error("Error reading body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	order := b.String()
	// todo валидация номера заказа
	// todo запись в базу заказа
	// todo асинхронное получение балов лояльности по заказу

	fmt.Fprint(w, order)
}
