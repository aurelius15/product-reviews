package rest

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/aurelius15/product-reviews/internal/service"
	"github.com/aurelius15/product-reviews/internal/web/rest/handler"
)

type Server struct {
	srv *http.Server
}

func StartRESTServer(ser *service.APIService) *Server {
	r := gin.New()

	apiGroup := r.Group("/api/v1")
	apiGroup.Use(gin.Logger())

	productAPIGroup := apiGroup.Group("/products")
	productHandler := handler.NewProductHandler(ser)

	productAPIGroup.GET("/:id", productHandler.Retrieve)
	productAPIGroup.POST("/:id", productHandler.Create)
	productAPIGroup.PUT("/:id", productHandler.Update)
	productAPIGroup.DELETE("/:id", productHandler.Delete)

	srv := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: time.Second,
		Handler:           r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error while starting server", slog.Any("err", err))
		}
	}()

	return &Server{srv: srv}
}

func (s Server) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("error while shutting down server", slog.Any("err", err))
	} else {
		slog.Info("server shut down")
	}
}
