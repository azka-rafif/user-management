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
	exists, err = s.ProductExistsInCart(cartId, prod.Id)
	if err != nil {
		return
	}
	if exists {
		res, err = s.UpdateCartItem(load, userId, cartId, prod.Id, prod.Price)
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
	items, err := s.Repo.GetCartItems(res.Id.String())
	if err != nil {
		return
	}
	res = res.AttachItems(items)
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
	orderItemsPayload := []order.OrderItemPayload{}
	var total float64
	for _, id := range load.CartItemsIds {
		exists, err = s.Repo.CartItemExistsByID(id)
		if err != nil {
			return
		}
		if !exists {
			err = failure.NotFound("Cart item")
			return
		}
		item, err := s.Repo.GetCartItemsByID(id)
		if err != nil {
			return res, err
		}
		exists, err = s.ProductService.ExistsByID(item.ProductId)
		if err != nil {
			return res, err
		}
		if !exists {
			err = failure.NotFound("Product")
			return res, err
		}
		orderItemsPayload = append(orderItemsPayload, order.OrderItemPayload{
			ProductId: item.ProductId,
			UserId:    userId,
			OrderId:   res.Id,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
		total += item.Price
	}
	res, err = s.OrderService.CreateOrder(order.OrderPayload{
		UserId:     userId,
		TotalPrice: total,
		Status:     "pending",
	}, orderItemsPayload)
	if err != nil {
		return
	}
	err = s.Repo.Checkout(load.CartItemsIds)
	return
}

func (s *CartServiceImpl) UpdateCartItem(load CartItemPayload, userId, cartId, productId uuid.UUID, productPrice float64) (res CartItem, err error) {
	res, err = s.Repo.GetCartItemByProduct(productId.String(), cartId.String())
	if err != nil {
		return
	}
	if (res.Quantity + load.Quantity) < 0 {
		err = failure.BadRequest(errors.New("quantity cannot be less than 0"))
		return
	}
	res.Update(load, userId)
	res.Recalculate(productPrice)
	err = s.Repo.UpdateItem(res)
	return
}

func (s *CartServiceImpl) ProductExistsInCart(cartId, prodId uuid.UUID) (exists bool, err error) {
	exists, err = s.Repo.ProductExistsInCart(prodId.String(), cartId.String())
	if err != nil {
		return
	}
	return
}
