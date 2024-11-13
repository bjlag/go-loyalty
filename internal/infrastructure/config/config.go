package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	defaultLogLevel          = "INFO"
	defaultRunAddrHost       = "localhost"
	defaultRunAddrPort       = 8080
	defaultJWTSecretKey      = "secret"
	defaultJWTExpTime        = 3 * time.Hour
	defaultDatabaseURI       = "postgres://postgres:secret@localhost:5432/master?sslmode=disable"
	defaultMigratePath       = "./migrations"
	defaultAccrualSystemHost = "localhost"
	defaultAccrualSystemPort = 9090

	envRunAddress           = "RUN_ADDRESS"
	envLogLevel             = "LOG_LEVEL"
	envJWTSecretKey         = "JWT_SECRET_KEY"
	envJWTExpTime           = "JWT_EXP_TIME"
	envDatabaseURI          = "DATABASE_URI"
	envMigratePath          = "MIGRATE_SOURCE_PATH"
	envAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
)

var (
	logLevel     string
	runAddr      *addr
	jwtSecretKey string
	jwtExpTime   time.Duration
	databaseURI  string
	migratePath  string
	accrualAddr  *addr
)

type Configuration struct {
	runAddr      addr
	logLevel     string
	jwtSecretKey string
	jwtExpTime   time.Duration
	databaseURI  string
	migratePath  string
	accrualAddr  addr
}

func Parse() *Configuration {
	runAddr = &addr{
		host: defaultRunAddrHost,
		port: defaultRunAddrPort,
	}
	accrualAddr = &addr{
		host: defaultAccrualSystemHost,
		port: defaultAccrualSystemPort,
	}
	jwtExpTime = defaultJWTExpTime

	parseFlags()
	parseEnvs()

	return &Configuration{
		runAddr:      *runAddr,
		logLevel:     logLevel,
		jwtSecretKey: jwtSecretKey,
		jwtExpTime:   jwtExpTime,
		databaseURI:  databaseURI,
		migratePath:  migratePath,
		accrualAddr:  *accrualAddr,
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

func (c Configuration) DatabaseURI() string {
	return c.databaseURI
}

func (c Configuration) MigratePath() string {
	return c.migratePath
}

func (c Configuration) AccrualSystemHost() string {
	return c.accrualAddr.host
}

func (c Configuration) AccrualSystemPort() int {
	return c.accrualAddr.port
}

func parseFlags() {
	var err error

	flag.StringVar(&logLevel, "l", defaultLogLevel, "Log level")
	flag.StringVar(&jwtSecretKey, "s", defaultJWTSecretKey, "JWT secret key")
	flag.StringVar(&databaseURI, "d", defaultDatabaseURI, "Database URI")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "Path to migration source files")
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
			if runAddr, err = newAddr(s); err != nil {
				return err
			}

			return nil
		},
	)
	flag.Func(
		"r",
		fmt.Sprintf("Accrual systm address: host:port (default \"%s:%d\")", defaultAccrualSystemHost, defaultAccrualSystemPort),
		func(s string) error {
			if accrualAddr, err = newAddr(s); err != nil {
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

	if value := os.Getenv(envDatabaseURI); value != "" {
		databaseURI = value
	}

	if value := os.Getenv(envMigratePath); value != "" {
		migratePath = value
	}

	if value := os.Getenv(envJWTExpTime); value != "" {
		if jwtExpTime, err = time.ParseDuration(value); err != nil {
			logEnvError(envJWTExpTime, value, err)
			os.Exit(2)
		}
	}

	if value := os.Getenv(envRunAddress); value != "" {
		if runAddr, err = newAddr(value); err != nil {
			logEnvError(envRunAddress, value, err)
			os.Exit(2)
		}
	}

	if value := os.Getenv(envAccrualSystemAddress); value != "" {
		if accrualAddr, err = newAddr(value); err != nil {
			logEnvError(envAccrualSystemAddress, value, err)
			os.Exit(2)
		}
	}
}

func logEnvError(env, value string, err error) {
	log.Fatalf("failed to parse environment variable %s=%s: %v", env, value, err)
}
