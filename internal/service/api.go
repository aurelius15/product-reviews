package service

import (
	"github.com/aurelius15/product-reviews/internal/storage"
)

type APIService struct {
	ProductService
	ReviewService
}

func NewAPIService(db *storage.PostgresStorage) *APIService {
	return &APIService{
		ProductService: NewProductService(db),
		ReviewService:  NewReviewService(db),
	}
}
