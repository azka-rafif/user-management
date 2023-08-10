package order

import "github.com/gofrs/uuid"

type OrderService interface {
	CreateOrder(load OrderPayload, itemLoads []OrderItemPayload) (res Order, err error)
	CreateOrderItem(load OrderItemPayload) (res OrderItem, err error)
	GetAll(limit, offset int, sort, field, status string, userId uuid.UUID, userRole string) (res []Order, err error)
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

func (s *OrderServiceImpl) GetAll(limit, offset int, sort, field, status string, userId uuid.UUID, userRole string) (res []Order, err error) {
	res, err = s.Repo.GetAll(limit, offset, sort, field, status, userId.String(), userRole)
	if err != nil {
		return
	}
	return
}
