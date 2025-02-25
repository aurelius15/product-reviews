package storage

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/aurelius15/product-reviews/internal/config"
)

type DataStore interface {
	Instance() *gorm.DB
	Close() error
}

type PostgresStorage struct {
	instance *gorm.DB
}

var instance *PostgresStorage

func NewPostgresStorage(cnf *config.PostgresCnf) (*PostgresStorage, error) {
	if instance != nil {
		return instance, nil
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cnf.ConnectionString(),
		PreferSimpleProtocol: true,
	}))

	if err != nil {
		return nil, err
	}

	instance = &PostgresStorage{
		instance: db,
	}

	slog.Info("db connection is created")

	return instance, nil
}

func (p PostgresStorage) Instance() *gorm.DB {
	return p.instance
}

func (p PostgresStorage) Close() error {
	if p.instance == nil {
		return nil
	}

	sqlDB, err := p.instance.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	slog.Info("db connection is closed")

	return nil
}
