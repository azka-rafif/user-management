package cart

import (
	"errors"

	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type CartService interface {
	AddToCart(load CartItemPayload, userId, cartId uuid.UUID) (res CartItem, err error)
	GetCart(cartId uuid.UUID) (res Cart, err error)
	Checkout(load CheckoutPayload, cartId, userId uuid.UUID) (res order.Order, err error)
}

type CartServiceImpl struct {
	Repo           CartRepository
	ProductService product.ProductService
	OrderService   order.OrderService
}

func ProvideCartServiceImpl(repo CartRepository, proService product.ProductService, ordService order.OrderService) *CartServiceImpl {
	return &CartServiceImpl{Repo: repo, ProductService: proService, OrderService: ordService}
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
	prod, err := s.ProductService.GetByID(load.ProductId)
	if err != nil {
		return
	}
	if prod.Stock < load.Quantity {
		err = failure.BadRequest(errors.New("not enough stock available"))
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

func (s *CartServiceImpl) Checkout(load CheckoutPayload, cartId, userId uuid.UUID) (res order.Order, err error) {
	exists, err := s.Repo.CartExistsByID(cartId.String())
	if err != nil {
		return
	}
	if !exists {
		err = failure.NotFound("Cart")
		return
	}
	for _, id := range load.CartItemsIds {
		exists, err = s.Repo.CartItemExistsByID(id)
		if err != nil {
			return
		}
		if !exists {
			err = failure.NotFound("Cart item")
			return
		}
	}
	return
}
