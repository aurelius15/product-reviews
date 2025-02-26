package service

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/aurelius15/product-reviews/internal/nats"
	"github.com/aurelius15/product-reviews/internal/repository"
	"github.com/aurelius15/product-reviews/internal/repository/model"
	"github.com/aurelius15/product-reviews/internal/storage"
	"github.com/aurelius15/product-reviews/internal/web/rest/apimodel"
)

type ReviewService interface {
	RetrieveReview(id int) (*apimodel.Review, error)
	RetrieveProductReviews(productID int) ([]*apimodel.Review, error)
	SaveReview(apiReview *apimodel.Review) (*apimodel.Review, error)
	DeleteReview(id int) error
}
type reviewService struct {
	reviewRepo repository.ReviewRepository
	cache      storage.CacheStore
	publisher  nats.Publisher
	validate   *validator.Validate
}

func NewReviewService(db storage.DataStore, cache storage.CacheStore, publisher nats.Publisher) ReviewService {
	return &reviewService{
		reviewRepo: repository.NewReviewRepository(db),
		cache:      cache,
		publisher:  publisher,
		validate:   validator.New(),
	}
}

func (s *reviewService) RetrieveReview(id int) (*apimodel.Review, error) {
	review, err := s.reviewRepo.Get(id)
	if err != nil {
		return nil, err
	}

	aReview := &apimodel.Review{
		ID:        review.ID,
		FirstName: review.FirstName,
		LastName:  review.LastName,
		Comment:   review.Comment,
		Rating:    review.Rating,
		ProductID: review.ProductID,
	}

	return aReview, nil
}

func (s *reviewService) RetrieveProductReviews(productID int) ([]*apimodel.Review, error) {
	cacheKey := fmt.Sprintf("product:%d:reviews", productID)
	apiReviews := make([]*apimodel.Review, 0)

	if err := s.cache.Get(cacheKey, &apiReviews); err == nil {
		return apiReviews, nil
	}

	return s.cacheProductReviews(productID, cacheKey)
}

func (s *reviewService) SaveReview(apiReview *apimodel.Review) (*apimodel.Review, error) {
	if err := s.validate.Struct(apiReview); err != nil {
		return nil, err
	}

	review := &model.Review{
		ID:        apiReview.ID,
		FirstName: apiReview.FirstName,
		LastName:  apiReview.LastName,
		Comment:   apiReview.Comment,
		Rating:    apiReview.Rating,
		ProductID: apiReview.ProductID,
	}

	var err error
	if review.ID == 0 {
		review, err = s.reviewRepo.Create(review)
	} else {
		review, err = s.reviewRepo.Update(review)
	}

	if err != nil {
		return nil, err
	}

	aReview := &apimodel.Review{
		ID:        review.ID,
		FirstName: review.FirstName,
		LastName:  review.LastName,
		Comment:   review.Comment,
		Rating:    review.Rating,
		ProductID: review.ProductID,
	}

	_ = s.publisher.Publish(aReview, false)

	return aReview, nil
}

func (s *reviewService) DeleteReview(id int) error {
	err := s.reviewRepo.Delete(id)
	if err != nil {
		return err
	}

	_ = s.publisher.Publish(&apimodel.Review{
		ID: id,
	}, true)

	return nil
}

func (s *reviewService) cacheProductReviews(productID int, cacheKey string) ([]*apimodel.Review, error) {
	acquired, err := s.cache.Lock(cacheKey, 10*time.Second)
	if err != nil {
		return nil, err
	}

	if !acquired {
		time.Sleep(100 * time.Millisecond)
		return s.RetrieveProductReviews(productID)
	}
	defer s.cache.Unlock(cacheKey)

	reviews, err := s.reviewRepo.GetByProduct(productID)
	if err != nil {
		return nil, err
	}

	apiReviews := make([]*apimodel.Review, 0, len(reviews))
	for _, review := range reviews {
		apiReviews = append(apiReviews, &apimodel.Review{
			ID:        review.ID,
			FirstName: review.FirstName,
			LastName:  review.LastName,
			Comment:   review.Comment,
			Rating:    review.Rating,
			ProductID: review.ProductID,
		})
	}

	if err := s.cache.Set(cacheKey, apiReviews, 5*time.Minute); err != nil {
		return nil, err
	}

	return apiReviews, nil
}
