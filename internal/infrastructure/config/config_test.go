package config_test

import (
	"flag"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bjlag/go-loyalty/internal/infrastructure/config"
)

func resetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	flag.CommandLine.Usage = func() { flag.Usage() }
	flag.Usage = usage
}

func TestParse_Default(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd"}

	got := config.Parse()

	assert.Equal(t, "localhost", got.RunAddrHost())
	assert.Equal(t, 8080, got.RunAddrPort())
	assert.Equal(t, "INFO", got.LogLevel())
	assert.Equal(t, "secret", got.JWTSecretKey())
	assert.Equal(t, 3*time.Hour, got.JWTExpTime())
	assert.Equal(t, "postgres://postgres:secret@localhost:5432/master?sslmode=disable", got.DatabaseURI())
	assert.Equal(t, "./migrations", got.MigratePath())
}

func TestParse_Flags(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:8888",
		"-l", "DEBUG",
		"-s", "new_secret",
		"-e", "1h",
		"-d", "new_db_uri",
		"-m", "new_migration_path",
	}

	got := config.Parse()

	assert.Equal(t, "127.0.0.1", got.RunAddrHost())
	assert.Equal(t, 8888, got.RunAddrPort())
	assert.Equal(t, "DEBUG", got.LogLevel())
	assert.Equal(t, "new_secret", got.JWTSecretKey())
	assert.Equal(t, 1*time.Hour, got.JWTExpTime())
	assert.Equal(t, "new_db_uri", got.DatabaseURI())
	assert.Equal(t, "new_migration_path", got.MigratePath())
}

func TestParse_Envs(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd"}

	envs := map[string]string{
		"RUN_ADDRESS":         "127.0.0.1:8888",
		"LOG_LEVEL":           "DEBUG",
		"JWT_SECRET_KEY":      "new_secret",
		"JWT_EXP_TIME":        "1h",
		"DATABASE_URI":        "new_db_uri",
		"MIGRATE_SOURCE_PATH": "new_migration_path",
	}

	for e, v := range envs {
		defer func() {
			_ = os.Unsetenv(e)
		}()
		err := os.Setenv(e, v)
		require.NoError(t, err)
	}

	got := config.Parse()

	assert.Equal(t, "127.0.0.1", got.RunAddrHost())
	assert.Equal(t, 8888, got.RunAddrPort())
	assert.Equal(t, "DEBUG", got.LogLevel())
	assert.Equal(t, "new_secret", got.JWTSecretKey())
	assert.Equal(t, 1*time.Hour, got.JWTExpTime())
	assert.Equal(t, "new_db_uri", got.DatabaseURI())
	assert.Equal(t, "new_migration_path", got.MigratePath())

}

func TestParse_EnvsOverwriteFlags(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:8888",
		"-l", "DEBUG",
		"-s", "new_secret",
		"-e", "1h",
		"-d", "new_db_uri",
		"-m", "new_migration_path",
	}

	envs := map[string]string{
		"RUN_ADDRESS":         "127.0.0.5:9999",
		"LOG_LEVEL":           "ERROR",
		"JWT_SECRET_KEY":      "new_secret_from_env",
		"JWT_EXP_TIME":        "5h",
		"DATABASE_URI":        "new_db_uri_from_env",
		"MIGRATE_SOURCE_PATH": "new_migration_path_from_env",
	}

	for e, v := range envs {
		defer func() {
			_ = os.Unsetenv(e)
		}()
		err := os.Setenv(e, v)
		require.NoError(t, err)
	}

	got := config.Parse()

	assert.Equal(t, "127.0.0.5", got.RunAddrHost())
	assert.Equal(t, 9999, got.RunAddrPort())
	assert.Equal(t, "ERROR", got.LogLevel())
	assert.Equal(t, "new_secret_from_env", got.JWTSecretKey())
	assert.Equal(t, 5*time.Hour, got.JWTExpTime())
	assert.Equal(t, "new_db_uri_from_env", got.DatabaseURI())
	assert.Equal(t, "new_migration_path_from_env", got.MigratePath())

}
