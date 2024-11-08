package login

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	errInvalidLogin    = errors.New("invalid login")
	errInvalidPassword = errors.New("invalid password")
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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
	if r.Login == "" {
		errs = append(errs, fmt.Errorf("%w: empty login", errInvalidLogin))
	}

	if r.Password == "" {
		errs = append(errs, fmt.Errorf("%w: empty password", errInvalidPassword))
	}

	return errors.Join(errs...)
}
