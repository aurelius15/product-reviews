package service

import (
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
}

func NewProductService(db *storage.PostgresStorage) ProductService {
	return &productService{
		productRepo: repository.NewProductRepository(db),
	}
}

func (s *productService) RetrieveProduct(id int) (*apimodel.Product, error) {
	p, err := s.productRepo.Get(id)
	if err != nil {
		return nil, err
	}

	aProduct := &apimodel.Product{
		ID:    p.ID,
		Name:  p.Name,
		Desc:  p.Description,
		Price: p.Price,
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
