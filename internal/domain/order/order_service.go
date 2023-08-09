package order

type OrderService interface {
	Create(load OrderPayload) (res Order, err error)
}

type OrderServiceImpl struct {
	Repo OrderRepository
}

func ProvideOrderServiceImpl(repo OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{Repo: repo}
}

func (s *OrderServiceImpl) Create(load OrderPayload) (res Order, err error) {
	return
}
