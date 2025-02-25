package migration

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/aurelius15/product-reviews/internal/config"
)

func RunMigrationCmd(postgresCnf *config.PostgresCnf) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", postgresCnf.PathToMigrations()),
		postgresCnf.ConnectionString(),
	)

	if err != nil {
		return err
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			slog.Error("closing migration process", slog.Any("err", err))
		}

		if dbErr != nil {
			slog.Error("closing migration process", slog.Any("err", err))
		}
	}()

	upMigrationErr := m.Up()
	if upMigrationErr == nil || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
