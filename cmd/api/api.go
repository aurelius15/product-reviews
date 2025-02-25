package api

import (
	"context"

	"github.com/aurelius15/product-reviews/internal/config"
	"github.com/aurelius15/product-reviews/internal/storage"
)

func RestAPICmd(_ context.Context, postgresCnf *config.PostgresCnf) error {
	_, err := storage.NewPostgresStorage(postgresCnf)
	if err != nil {
		return err
	}

	return nil
}
