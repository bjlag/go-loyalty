package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

const (
	// database requirement
	maxLenEmail = 50
	// bcrypt: password length should not exceed 72 bytes
	maxLenPassword = 72
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

	if len(r.Login) > maxLenEmail {
		errs = append(errs, fmt.Errorf("%w: email length exceeds 50 bytes", errInvalidLogin))
	}

	if r.Password == "" {
		errs = append(errs, fmt.Errorf("%w: empty password", errInvalidPassword))
	}

	if len(r.Password) > maxLenPassword {
		errs = append(errs, fmt.Errorf("%w: password length exceeds 72 bytes", errInvalidPassword))
	}

	return errors.Join(errs...)
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
