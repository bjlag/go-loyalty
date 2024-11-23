package withdraw

import (
	"encoding/json"
	"errors"

	"github.com/bjlag/go-loyalty/internal/infrastructure/validator"
)

var (
	errInvalidOrder = errors.New("invalid order")
	errInvalidSum   = errors.New("invalid sum")
)

type Request struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (r *Request) UnmarshalJSON(b []byte) error {
	type RequestAlias Request

	aliasValue := &struct {
		*RequestAlias
	}{
		RequestAlias: (*RequestAlias)(r),
	}

	err := json.Unmarshal(b, &aliasValue)
	if err != nil {
		return err
	}

	var errs []error

	if !validator.CheckLuhn(r.Order) {
		errs = append(errs, errInvalidOrder)
	}

	if r.Sum <= 0 {
		errs = append(errs, errInvalidSum)
	}

	return errors.Join(errs...)
}
