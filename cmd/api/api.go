package api

import (
	"context"

	"github.com/aurelius15/product-reviews/internal/config"
	"github.com/aurelius15/product-reviews/internal/service"
	"github.com/aurelius15/product-reviews/internal/storage"
	"github.com/aurelius15/product-reviews/internal/utils"
	"github.com/aurelius15/product-reviews/internal/web/rest"
)

func RestAPICmd(ctx context.Context, postgresCnf *config.PostgresCnf) error {
	postgresStorage, err := storage.NewPostgresStorage(postgresCnf)
	if err != nil {
		return err
	}

	s := rest.StartRESTServer(service.NewAPIService(postgresStorage))

	gCtx, cancel := utils.GracefulShutdown(ctx)
	defer cancel()

	s.Shutdown(gCtx)
	postgresStorage.Close()

	return nil
}
