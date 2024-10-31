package config

import (
	"flag"
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

type Configuration struct {
	runAddr  runAddr
	logLevel string
}

func Parse() *Configuration {
	var (
		err      error
		logLevel string
	)

	addr := &runAddr{
		host: defaultRunAddrHost,
		port: defaultRunAddrPort,
	}

	flag.Func("a", "Server address: host:port", func(s string) error {
		if addr, err = newRunAddr(s); err != nil {
			return err
		}

		return nil
	})

	flag.StringVar(&logLevel, "l", defaultLogLevel, "Log level")

	flag.Parse()

	if envRunAddrValue := os.Getenv(envRunAddress); envRunAddrValue != "" {
		if addr, err = newRunAddr(envRunAddrValue); err != nil {
			logEnvError(envRunAddress, envRunAddrValue, err)
			return nil
		}
	}

	if envLogLevelValue := os.Getenv(envLogLevel); envLogLevelValue != "" {
		logLevel = envLogLevelValue
	}

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

func logEnvError(env, value string, err error) {
	log.Fatalf("failed to parse environment variable %s=%s: %v", env, value, err)
}
