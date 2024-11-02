package migrator

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(sourcePath string, databaseInstance database.Driver) (*Migrator, error) {
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", sourcePath),
		"master",
		databaseInstance,
	)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		migrate: m,
	}, nil
}

func (m Migrator) Up() (bool, error) {
	if err := m.migrate.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
