package cart

import (
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type CartService interface {
	AddToCart(load CartItemPayload, userId, cartId uuid.UUID) (res CartItem, err error)
	GetCart(cartId uuid.UUID) (res Cart, err error)
}

type CartServiceImpl struct {
	Repo           CartRepository
	ProductService product.ProductService
}

func ProvideCartServiceImpl(repo CartRepository, proService product.ProductService) *CartServiceImpl {
	return &CartServiceImpl{Repo: repo, ProductService: proService}
}

func (s *CartServiceImpl) AddToCart(load CartItemPayload, userId, cartId uuid.UUID) (res CartItem, err error) {
	exists, err := s.Repo.CartExistsByID(cartId.String())
	if err != nil {
		return
	}
	if !exists {
		err = failure.NotFound("Cart")
		return
	}

	prodId, err := uuid.FromString(load.ProductId.String())
	if err != nil {
		return
	}
	prod, err := s.ProductService.GetByID(prodId)
	if err != nil {
		return
	}
	res, err = res.NewFromPayload(load, cartId, userId, prod.Price)
	if err != nil {
		return
	}
	err = s.Repo.CreateItem(res)
	return
}

func (s *CartServiceImpl) GetCart(cartId uuid.UUID) (res Cart, err error) {
	exists, err := s.Repo.CartExistsByID(cartId.String())

	if err != nil {
		return
	}

	if !exists {
		err = failure.NotFound("Cart")
		return
	}
	res, err = s.Repo.GetCartByID(cartId.String())
	if err != nil {
		return
	}
	return
}

func (s *CartServiceImpl) GetCartItems(cartId uuid.UUID) (res []CartItem, err error) {
	exists, err := s.Repo.CartExistsByID(cartId.String())
	if err != nil {
		return
	}

	if !exists {
		err = failure.NotFound("Cart")
		return
	}
	res, err = s.Repo.GetCartItems(cartId.String())
	if err != nil {
		return
	}
	return
}
