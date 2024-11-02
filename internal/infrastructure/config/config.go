package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	defaultRunAddrHost  = "localhost"
	defaultRunAddrPort  = 8080
	defaultLogLevel     = "INFO"
	defaultJWTSecretKey = "secret"
	defaultJWTExpTime   = 3 * time.Hour

	envRunAddress   = "RUN_ADDRESS"
	envLogLevel     = "LOG_LEVEL"
	envJWTSecretKey = "JWT_SECRET_KEY"
	envJWTExpTime   = "JWT_EXP_TIME"
)

var (
	logLevel     string
	addr         *runAddr
	jwtSecretKey string
	jwtExpTime   time.Duration
)

type Configuration struct {
	runAddr      runAddr
	logLevel     string
	jwtSecretKey string
	jwtExpTime   time.Duration
}

func Parse() *Configuration {
	addr = &runAddr{
		host: defaultRunAddrHost,
		port: defaultRunAddrPort,
	}
	jwtExpTime = defaultJWTExpTime

	parseFlags()
	parseEnvs()

	return &Configuration{
		runAddr:      *addr,
		logLevel:     logLevel,
		jwtSecretKey: jwtSecretKey,
		jwtExpTime:   jwtExpTime,
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

func (c Configuration) JWTSecretKey() string {
	return c.jwtSecretKey
}

func (c Configuration) JWTExpTime() time.Duration {
	return c.jwtExpTime
}

func parseFlags() {
	var err error

	flag.StringVar(&logLevel, "l", defaultLogLevel, "Log level")
	flag.StringVar(&jwtSecretKey, "s", defaultJWTSecretKey, "JWT secret key")
	flag.Func("e", "JWT token expiration time (default 3h)", func(s string) error {
		jwtExpTime, err = time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid JWT token expiration time: %w", err)
		}

		return nil
	})
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

	if value := os.Getenv(envJWTSecretKey); value != "" {
		jwtSecretKey = value
	}

	if value := os.Getenv(envJWTExpTime); value != "" {
		if jwtExpTime, err = time.ParseDuration(value); err != nil {
			logEnvError(envJWTExpTime, value, err)
			os.Exit(2)
		}

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
