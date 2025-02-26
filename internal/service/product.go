package service

import (
	"fmt"
	"time"

	"github.com/aurelius15/product-reviews/internal/repository"
	"github.com/aurelius15/product-reviews/internal/repository/model"
	"github.com/aurelius15/product-reviews/internal/storage"
	"github.com/aurelius15/product-reviews/internal/web/rest/apimodel"
)

type ProductService interface {
	RetrieveProduct(id int) (*apimodel.Product, error)
	SaveProduct(apiProduct *apimodel.Product) (*apimodel.Product, error)
	DeleteProduct(id int) error
}

type productService struct {
	productRepo repository.ProductRepository
	cache       storage.CacheStore
}

func NewProductService(db storage.DataStore, cache storage.CacheStore) ProductService {
	return &productService{
		productRepo: repository.NewProductRepository(db),
		cache:       cache,
	}
}

func (s *productService) RetrieveProduct(id int) (*apimodel.Product, error) {
	p, err := s.productRepo.Get(id)
	if err != nil {
		return nil, err
	}

	avg, err := s.calculateAvgRating(p.ID)
	if err != nil {
		return nil, err
	}

	aProduct := &apimodel.Product{
		ID:        p.ID,
		Name:      p.Name,
		Desc:      p.Description,
		Price:     p.Price,
		AvgRating: avg,
	}

	return aProduct, nil
}

func (s *productService) SaveProduct(apiProduct *apimodel.Product) (*apimodel.Product, error) {
	product := &model.Product{
		ID:          apiProduct.ID,
		Name:        apiProduct.Name,
		Description: apiProduct.Desc,
		Price:       apiProduct.Price,
	}

	var err error
	if product.ID == 0 {
		product, err = s.productRepo.Create(product)
	} else {
		product, err = s.productRepo.Update(product)
	}

	if err != nil {
		return nil, err
	}

	return &apimodel.Product{
		ID:    product.ID,
		Name:  product.Name,
		Desc:  product.Description,
		Price: product.Price,
	}, nil
}

func (s *productService) DeleteProduct(id int) error {
	err := s.productRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// calculateAvgRating retrieves the average rating for the specified product ID from the cache or calculates and caches it.
// Returns the average rating or an error if the operation fails.
func (s *productService) calculateAvgRating(productID int) (float64, error) {
	cacheKey := fmt.Sprintf("product:%d:rating", productID)

	var avg float64

	if err := s.cache.Get(cacheKey, &avg); err == nil {
		return avg, nil
	}

	return s.recalculateAvgRatingAndCache(productID, cacheKey)
}

// recalculateAvgRatingAndCache recalculates and caches the average rating for the given product ID with a specified cache key.
// It acquires a lock to prevent concurrent modifications and retrieves the rating from the cache or repository as needed.
// Returns the calculated average rating or an error if the operation fails.
func (s *productService) recalculateAvgRatingAndCache(productID int, cacheKey string) (float64, error) {
	acquired, err := s.cache.Lock(cacheKey, 10*time.Second)
	if err != nil {
		return 0, err
	}

	if !acquired {
		time.Sleep(100 * time.Millisecond)
		return s.calculateAvgRating(productID)
	}
	defer s.cache.Unlock(cacheKey)

	var avg float64
	if err := s.cache.Get(cacheKey, &avg); err == nil {
		return avg, nil
	}

	avg, err = s.productRepo.GetAvgRating(productID)
	if err != nil {
		return 0, err
	}

	if err := s.cache.Set(cacheKey, avg, 5*time.Minute); err != nil {
		return 0, err
	}

	return avg, nil
}
