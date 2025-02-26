package repository

import (
	"github.com/pkg/errors"

	"github.com/aurelius15/product-reviews/internal/repository/model"
	"github.com/aurelius15/product-reviews/internal/storage"
)

type ReviewRepository interface {
	Create(p *model.Review) (*model.Review, error)
	Get(id int) (*model.Review, error)
	GetByProduct(productID int) ([]*model.Review, error)
	Update(p *model.Review) (*model.Review, error)
	Delete(id int) error
}
type reviewRepository struct {
	db storage.DataStore
}

func NewReviewRepository(db storage.DataStore) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *model.Review) (*model.Review, error) {
	if review == nil {
		return nil, errors.New("review is nil")
	}

	if review.ID != 0 {
		return nil, errors.New("can't create review with id")
	}

	result := r.db.Instance().Create(review)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to review product")
	}

	return review, nil
}

func (r *reviewRepository) Get(id int) (*model.Review, error) {
	review := &model.Review{}
	result := r.db.Instance().First(review, id)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to get review")
	}

	return review, nil
}

func (r *reviewRepository) GetByProduct(productID int) ([]*model.Review, error) {
	var reviews []*model.Review
	result := r.db.Instance().Find(&reviews, "product_id = ?", productID)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to get reviews")
	}

	return reviews, nil
}

func (r *reviewRepository) Update(review *model.Review) (*model.Review, error) {
	if review == nil {
		return nil, errors.New("review is nil")
	}

	id := review.ID
	review.ID = 0

	if id == 0 {
		return nil, errors.New("can't update review without id")
	}

	result := r.db.Instance().Model(review).Where("id = ?", id).Updates(review)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to update review")
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("failed to update review")
	}

	review.ID = id

	return review, nil
}

func (r *reviewRepository) Delete(id int) error {
	result := r.db.Instance().Delete(&model.Review{}, id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete review")
	}

	return nil
}
