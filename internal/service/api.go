package service

import (
	"github.com/aurelius15/product-reviews/internal/nats"
	"github.com/aurelius15/product-reviews/internal/storage"
)

type APIService struct {
	ProductService
	ReviewService
}

func NewAPIService(db storage.DataStore, publisher nats.Publisher) *APIService {
	return &APIService{
		ProductService: NewProductService(db),
		ReviewService:  NewReviewService(db, publisher),
	}
}
