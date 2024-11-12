package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type addr struct {
	host string
	port int
}

func newAddr(s string) (*addr, error) {
	values := strings.Split(s, ":")
	if len(values) != 2 {
		return nil, errors.New("invalid format")
	}

	port, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	return &addr{
		host: values[0],
		port: port,
	}, nil
}

func (a addr) Host() string {
	return a.host
}

func (a addr) Port() int {
	return a.port
}
