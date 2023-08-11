package order

import (
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type OrderService interface {
	CreateOrder(load OrderPayload, itemLoads []OrderItemPayload) (res Order, err error)
	CreateOrderItem(load OrderItemPayload) (res OrderItem, err error)
	GetAll(limit, offset int, sort, field, status string, userId uuid.UUID, userRole string, cancelled bool) (res []Order, err error)
	CancelOrder(orderId, userId uuid.UUID, userRole string) (res Order, err error)
	GetByID(orderId, userId uuid.UUID) (res Order, err error)
}

type OrderServiceImpl struct {
	Repo OrderRepository
}

func ProvideOrderServiceImpl(repo OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{Repo: repo}
}

func (s *OrderServiceImpl) CreateOrder(load OrderPayload, itemLoads []OrderItemPayload) (res Order, err error) {
	res, err = res.NewFromPayload(load)
	if err != nil {
		return
	}
	items := []OrderItem{}
	for _, itemLoad := range itemLoads {
		itemLoad.OrderId = res.Id
		item, err := s.CreateOrderItem(itemLoad)
		if err != nil {
			return res, err
		}
		items = append(items, item)
	}
	res = res.AttachItems(items)
	err = s.Repo.Create(res)

	return
}

func (s *OrderServiceImpl) CreateOrderItem(load OrderItemPayload) (res OrderItem, err error) {
	res, err = res.NewFromPayload(load)
	if err != nil {
		return
	}
	return
}

func (s *OrderServiceImpl) GetAll(limit, offset int, sort, field, status string, userId uuid.UUID, userRole string, cancelled bool) (res []Order, err error) {
	res, err = s.Repo.GetAll(limit, offset, sort, field, status, userId.String(), userRole, cancelled)
	if err != nil {
		return
	}

	for i, order := range res {
		items, err := s.Repo.GetItemsByOrderID(order.Id.String())
		if err != nil {
			return res, err
		}
		order = order.AttachItems(items)
		res[i] = order
	}

	return
}

func (s *OrderServiceImpl) CancelOrder(orderId, userId uuid.UUID, userRole string) (res Order, err error) {
	exists, err := s.Repo.ExistsByID(orderId.String())
	if err != nil {
		return
	}
	if !exists {
		err = failure.NotFound("Order")
		return
	}
	res, err = s.Repo.GetOrderByID(orderId.String())
	if err != nil {
		return
	}
	if userId.String() != res.UserId.String() && userRole != "admin" {
		err = failure.Unauthorized("unauthorized, invalid credentials")
		return
	}
	items, err := s.Repo.GetItemsByOrderID(orderId.String())
	if err != nil {
		return
	}
	res = res.AttachItems(items)
	err = res.CancelOrder(userId)

	if err != nil {
		return
	}
	err = s.Repo.CancelOrder(res)
	return
}

func (s *OrderServiceImpl) GetByID(orderId, userId uuid.UUID) (res Order, err error) {
	res, err = s.Repo.GetOrderByID(orderId.String())
	if err != nil {
		return
	}
	if res.UserId != userId {
		err = failure.Unauthorized("Invalid Credentials")
		return
	}
	items, err := s.Repo.GetItemsByOrderID(orderId.String())
	if err != nil {
		return
	}
	res = res.AttachItems(items)
	return
}
