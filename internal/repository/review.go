package repository

import (
	"github.com/pkg/errors"

	"github.com/aurelius15/product-reviews/internal/repository/model"
	"github.com/aurelius15/product-reviews/internal/storage"
)

type ReviewRepository interface {
	Create(p *model.Review) error
	Get(id int) (*model.Review, error)
	Update(p *model.Review) error
	Delete(id int) error
}
type reviewRepository struct {
	db storage.DataStore
}

func NewReviewRepository(db storage.DataStore) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(p *model.Review) error {
	if p == nil {
		return errors.New("review is nil")
	}

	if p.ID != 0 {
		return errors.New("can't create review with id")
	}

	result := r.db.Instance().Create(p)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to create review")
	}

	return nil
}

func (r *reviewRepository) Get(id int) (*model.Review, error) {
	p := &model.Review{}
	result := r.db.Instance().First(p, id)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to get review")
	}

	return p, nil
}

func (r *reviewRepository) Update(p *model.Review) error {
	if p == nil {
		return errors.New("review is nil")
	}

	result := r.db.Instance().Save(p)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update review")
	}

	return nil
}

func (r *reviewRepository) Delete(id int) error {
	result := r.db.Instance().Delete(&model.Review{}, id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete review")
	}

	return nil
}
