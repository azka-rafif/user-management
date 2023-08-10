package product

import (
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type ProductService interface {
	Create(load ProductPayload, userId uuid.UUID) (res Product, err error)
	GetAll(limit, offset int, sort, field, productTitle string) (res []Product, err error)
	GetByID(id uuid.UUID) (res Product, err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
}

type ProductServiceImpl struct {
	Repo ProductRepository
}

func ProvideProductServiceImpl(repo ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{Repo: repo}
}

func (s *ProductServiceImpl) Create(load ProductPayload, userId uuid.UUID) (res Product, err error) {

	res, err = res.NewFromPayload(load, userId)
	if err != nil {
		return
	}
	err = s.Repo.Create(res)
	if err != nil {
		return
	}
	return
}

func (s *ProductServiceImpl) GetAll(limit, offset int, sort, field, productTitle string) (res []Product, err error) {
	res, err = s.Repo.GetAll(limit, offset, sort, field, productTitle)
	if err != nil {
		return
	}
	return
}

func (s *ProductServiceImpl) GetByID(id uuid.UUID) (res Product, err error) {
	exists, err := s.Repo.ExistsByID(id.String())

	if err != nil {
		return
	}

	if !exists {
		err = failure.NotFound("Product")
		return
	}

	res, err = s.Repo.GetByID(id.String())
	if err != nil {
		return
	}

	return
}

func (s *ProductServiceImpl) ExistsByID(id uuid.UUID) (exists bool, err error) {
	exists, err = s.Repo.ExistsByID(id.String())
	return
}
