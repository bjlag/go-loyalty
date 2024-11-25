package register

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// database requirement
	maxLenLogin = 20
	// bcrypt: password length should not exceed 72 bytes
	maxLenPassword = 72
	minLenPassword = 6
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

	if len(r.Login) > maxLenLogin {
		errs = append(errs, fmt.Errorf("%w: email length exceeds 50 bytes", errInvalidLogin))
	}

	if len(r.Password) < minLenPassword {
		errs = append(errs, fmt.Errorf("%w: password length less than %d bytes", errInvalidPassword, minLenPassword))
	}

	if len(r.Password) > maxLenPassword {
		errs = append(errs, fmt.Errorf("%w: password length exceeds %d bytes", errInvalidPassword, maxLenPassword))
	}

	return errors.Join(errs...)
}
