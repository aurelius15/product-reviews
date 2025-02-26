package repository

import (
	"github.com/pkg/errors"

	"github.com/aurelius15/product-reviews/internal/repository/model"
	"github.com/aurelius15/product-reviews/internal/storage"
)

type ProductRepository interface {
	Create(p *model.Product) (*model.Product, error)
	Get(id int) (*model.Product, error)
	Update(p *model.Product) (*model.Product, error)
	Delete(id int) error
}
type productRepository struct {
	db storage.DataStore
}

func NewProductRepository(db storage.DataStore) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *model.Product) (*model.Product, error) {
	if p == nil {
		return nil, errors.New("product is nil")
	}

	if p.ID != 0 {
		return nil, errors.New("can't create product with id")
	}

	result := r.db.Instance().Debug().Create(p)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to create product")
	}

	return p, nil
}

func (r *productRepository) Get(id int) (*model.Product, error) {
	p := &model.Product{}
	result := r.db.Instance().First(p, id)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to get product")
	}

	return p, nil
}

func (r *productRepository) Update(p *model.Product) (*model.Product, error) {
	if p == nil {
		return nil, errors.New("product is nil")
	}

	id := p.ID
	p.ID = 0

	if id == 0 {
		return nil, errors.New("can't update product without id")
	}

	result := r.db.Instance().Debug().Model(p).Where("id = ?", id).Updates(p)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to update product")
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("failed to update product")
	}

	p.ID = id

	return p, nil
}

func (r *productRepository) Delete(id int) error {
	result := r.db.Instance().Delete(&model.Product{}, id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete product")
	}

	return nil
}
