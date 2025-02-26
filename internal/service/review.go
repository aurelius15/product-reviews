package service

import (
	"github.com/aurelius15/product-reviews/internal/repository"
	"github.com/aurelius15/product-reviews/internal/repository/model"
	"github.com/aurelius15/product-reviews/internal/storage"
	"github.com/aurelius15/product-reviews/internal/web/rest/apimodel"
)

type ReviewService interface {
	RetrieveReview(id int) (*apimodel.Review, error)
	SaveReview(apiReview *apimodel.Review) (*apimodel.Review, error)
	DeleteReview(id int) error
}
type reviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(db *storage.PostgresStorage) ReviewService {
	return &reviewService{
		reviewRepo: repository.NewReviewRepository(db),
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

func (s *reviewService) SaveReview(apiReview *apimodel.Review) (*apimodel.Review, error) {
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

	return &apimodel.Review{
		ID:        review.ID,
		FirstName: review.FirstName,
		LastName:  review.LastName,
		Comment:   review.Comment,
		Rating:    review.Rating,
		ProductID: review.ProductID,
	}, nil
}

func (s *reviewService) DeleteReview(id int) error {
	err := s.reviewRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
