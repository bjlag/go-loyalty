package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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
	if !isEmailValid(r.Login) {
		errs = append(errs, errInvalidLogin)
	}

	if r.Password == "" {
		errs = append(errs, fmt.Errorf("%w: empty password", errInvalidPassword))
	}

	return errors.Join(errs...)
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
