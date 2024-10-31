package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	envRunAddress = "RUN_ADDRESS"
)

type RunAddr struct {
	host string
	port int
}

func newRunAddr(s string) (*RunAddr, error) {
	values := strings.Split(s, ":")
	if len(values) != 2 {
		return nil, errors.New("invalid format")
	}

	port, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	return &RunAddr{
		host: values[0],
		port: port,
	}, nil
}

func (a RunAddr) Host() string {
	return a.host
}

func (a RunAddr) Port() int {
	return a.port
}

type Configuration struct {
	runAddr RunAddr
}

func Parse() *Configuration {
	var (
		err  error
		addr *RunAddr
	)

	flag.Func("a", "Server address: host:port", func(s string) error {
		if addr, err = newRunAddr(s); err != nil {
			return err
		}

		return nil
	})

	flag.Parse()

	if envRunAddrValue := os.Getenv(envRunAddress); envRunAddrValue != "" {
		if addr, err = newRunAddr(envRunAddrValue); err != nil {
			logEnvError(envRunAddress, envRunAddrValue, err)
			return nil
		}
	}

	return &Configuration{
		runAddr: *addr,
	}
}

func (c Configuration) RunAddr() RunAddr {
	return c.runAddr
}

func logEnvError(env, value string, err error) {
	log.Fatalf("failed to parse environment variable %s=%s: %v", env, value, err)
}
