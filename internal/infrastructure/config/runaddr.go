package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type runAddr struct {
	host string
	port int
}

func newRunAddr(s string) (*runAddr, error) {
	values := strings.Split(s, ":")
	if len(values) != 2 {
		return nil, errors.New("invalid format")
	}

	port, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	return &runAddr{
		host: values[0],
		port: port,
	}, nil
}

func (a runAddr) Host() string {
	return a.host
}

func (a runAddr) Port() int {
	return a.port
}
