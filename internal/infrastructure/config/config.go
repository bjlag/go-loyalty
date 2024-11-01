package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	defaultRunAddrHost = "localhost"
	defaultRunAddrPort = 8080
	defaultLogLevel    = "INFO"

	envRunAddress = "RUN_ADDRESS"
	envLogLevel   = "LOG_LEVEL"
)

var (
	logLevel string
	addr     *runAddr
)

type Configuration struct {
	runAddr  runAddr
	logLevel string
}

func Parse() *Configuration {
	addr = &runAddr{
		host: defaultRunAddrHost,
		port: defaultRunAddrPort,
	}

	parseFlags()
	parseEnvs()

	return &Configuration{
		runAddr:  *addr,
		logLevel: logLevel,
	}
}

func (c Configuration) RunAddrHost() string {
	return c.runAddr.host
}

func (c Configuration) RunAddrPort() int {
	return c.runAddr.port
}

func (c Configuration) LogLevel() string {
	return c.logLevel
}

func parseFlags() {
	var err error

	flag.StringVar(&logLevel, "l", defaultLogLevel, "Log level")
	flag.Func(
		"a",
		fmt.Sprintf("Server address: host:port (default \"%s:%d\")", defaultRunAddrHost, defaultRunAddrPort),
		func(s string) error {
			if addr, err = newRunAddr(s); err != nil {
				return err
			}

			return nil
		},
	)

	flag.Parse()
}

func parseEnvs() {
	var err error

	if value := os.Getenv(envLogLevel); value != "" {
		logLevel = value
	}

	if value := os.Getenv(envRunAddress); value != "" {
		if addr, err = newRunAddr(value); err != nil {
			logEnvError(envRunAddress, value, err)
			os.Exit(2)
		}
	}
}

func logEnvError(env, value string, err error) {
	log.Fatalf("failed to parse environment variable %s=%s: %v", env, value, err)
}
