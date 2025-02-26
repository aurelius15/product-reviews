package storage

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/aurelius15/product-reviews/internal/config"
)

type DataStore interface {
	Instance() *gorm.DB
	Close()
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

func (p PostgresStorage) Close() {
	if p.instance == nil {
		return
	}

	sqlDB, err := p.instance.DB()
	if err != nil {
		slog.Error("db connection is not created")
		return
	}

	if err := sqlDB.Close(); err != nil {
		slog.Error("db connection is not closed", slog.Any("err", err))
		return
	}

	slog.Info("db connection is closed")
}
